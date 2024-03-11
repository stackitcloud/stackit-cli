{
  description = "Nix flake for the STACKIT CLI";

  inputs = {
    nixpkgs.url = "nixpkgs/nixpkgs-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      pname = "stackit-cli";

      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];

      # Helper functions
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {
      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          # The default package for 'nix run', 'nix profile install', 'nix build' & 'nix shell'
          default = pkgs.buildGoModule rec {
            # underlying 'stdenv.mkDerivation' requires 'name', if 'name' is not provided 'pname' & 'version' must exist!
            inherit pname self;
            version = "${self.rev or self.dirtyRev}";

            src = ./.; # TODO: check if I can change to a different revision here. This would reduce the maintenance pain of always needing to update `vendorHash`.
            # 'vendorHash' represents a derivative of all go.mod dependencies and needs to be adjusted with every change
            vendorHash = "sha256-2DJ/NBeYxAL501Nz4qru77hXtOhRCO9i6RRymPviYBg=";

            CGO_ENABLED = 0;
            ldflags = [ "-X main.version=${version}" ];

            subPackages = [ "/" ];

            nativeBuildInputs = [ pkgs.installShellFiles pkgs.makeWrapper ];

            # doCheck = false; # TODO: consider disabling checking to speed up building by ~10s
            # nativeCheckInputs = [ less ]; # TODO: check if I can really drop this
            preCheck = ''
              export HOME=$TMPDIR # TODO: needed because tests executes mkdir
            '';

            postInstall = ''
              export HOME=$TMPDIR # TODO: all invocations of the binary try to create a directory and file below $HOME
              mv $out/bin/{${pname},stackit} # rename the binary

              installShellCompletion --cmd stackit --bash <($out/bin/stackit completion bash)
              installShellCompletion --cmd stackit --zsh <($out/bin/stackit completion zsh)
              installShellCompletion --cmd stackit --fish <($out/bin/stackit completion fish)
              # Use this instead once LINK is fixed:
              # installShellCompletion --cmd stackit \
              #   --bash <($out/bin/stackit completion bash) \
              #   --zsh  <($out/bin/stackit completion zsh)  \
              #   --fish <($out/bin/stackit completion fish)
            '';

            # TODO: consider also wrapping `xdg-open` (needed for browser promped user login)
            postFixup = ''
              wrapProgram $out/bin/stackit \
                --prefix PATH : ${pkgs.lib.makeBinPath [ pkgs.less ]}
            '';

            # TODO: check if this can be run using flakes
            # passthru.tests = {
            #   simple = runCommand "${pname}-test" {} ''
            #     [ $(find ${stackit-cli}/share -not -empty -type f | wc -l) -eq 3 ]
            #   '';
            # };

            meta = with pkgs.lib; {
              description = "STACKIT CLI";
              homepage = "https://github.com/stackitcloud/stackit-cli";
              license = licenses.asl20;
              platforms = supportedSystems;
              mainProgram = "stackit";
            };
          };
        });
      overlays.default = final: prev: {
        ${pname} = self.packages.${prev.system}.default;
      };

      # 'nix develop'
      devShells = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [
              go
              gopls
              gotools
              go-tools
              gnumake
            ];
          };
        });

      # 'nix fmt'
      formatter = forAllSystems (system: nixpkgs.legacyPackages.${system}.nixpkgs-fmt);
    };
}
