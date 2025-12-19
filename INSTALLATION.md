# Installation

## Package managers

[![Packaging status](https://repology.org/badge/vertical-allrepos/stackit-cli.svg?columns=1)](https://repology.org/project/stackit-cli/versions)

### macOS

The STACKIT CLI can be installed through the [Homebrew](https://brew.sh/) package manager.

1. First, you need to register the [STACKIT tap](https://github.com/stackitcloud/homebrew-tap) via:

```shell
brew tap stackitcloud/tap
```

2. You can then install the CLI via:

```shell
brew install --cask stackit
```

#### Formula deprecated

The homebrew formula is deprecated, will no longer be updated and will be removed after 2026-01-22.
You need to install the STACKIT CLI as cask.  
Therefor you need to uninstall the formula and reinstall it as cask.  

Your profiles should normally remain. To ensure that nothing will be gone, you should backup them.

1. Export your existing profiles. This will create a json file in your current directory.
```shell
stackit config profile export default
```

2. If you have multiple profiles, then execute the export command for each of them. You can find your profiles via:

```shell
stackit config profile list
stackit config profile export <profile-name>
```

3. Uninstall the formula.
```shell
brew uninstall stackit
```

4. Install the STACKIT CLI as cask.
```shell
brew install --cask stackit
```

5. Check if your configs are still stored.
```shell
stackit config profile list
```

6. In case the profiles are gone, import your profiles via:
```shell
$ stackit config profile import -c @default.json --name myProfile
```

### Linux

#### Snapcraft

The STACKIT CLI is available as a [Snap](https://snapcraft.io/stackit), and can be installed via:

```shell
sudo snap install stackit --classic
```

or via the [Snap Store](https://snapcraft.io/snap-store) for desktop.

#### Debian/Ubuntu (`APT`)

The STACKIT CLI can be installed through the [`APT`](https://ubuntu.com/server/docs/package-management) package manager.

##### Before you begin

To install the STACKIT CLI package, you will need to have the `curl` and `gnupg` packages installed:

```shell
sudo apt-get update
sudo apt-get install curl gnupg
```

##### Installing

1. Import the STACKIT public key:

```shell
curl https://packages.stackit.cloud/keys/key.gpg | sudo gpg --dearmor -o /usr/share/keyrings/stackit.gpg
```

2. Add the STACKIT CLI package repository as a package source:

```shell
echo "deb [signed-by=/usr/share/keyrings/stackit.gpg] https://packages.stackit.cloud/apt/cli stackit main" | sudo tee -a /etc/apt/sources.list.d/stackit.list
```

3. Update repository information and install the `stackit` package:

```shell
sudo apt-get update
sudo apt-get install stackit
```

> If you can't install the `stackit` package due to an expired key, please go back to step `1` to import the latest public key.

#### Nix / NixOS

The STACKIT CLI is available as a [Nix package](https://search.nixos.org/packages?channel=unstable&show=stackit-cli), and can be used via:

```shell
nix-shell -p stackit-cli
```

#### Eget

The STACKIT CLI binaries are available via our [GitHub releases](https://github.com/stackitcloud/stackit-cli/releases), you can install them from there using [Eget](https://github.com/zyedidia/eget).

```toml
# ~/.eget.toml
["stackitcloud/stackit-cli"]
asset_filters=["stackit-cli_", "_linux_amd64.tar.gz"]
```

```shell
eget stackitcloud/stackit-cli
```

#### RHEL/Fedora/Rocky/Alma/openSUSE/... (`DNF/YUM/Zypper`)

The STACKIT CLI can be installed through the [`DNF/YUM`](https://docs.fedoraproject.org/en-US/fedora/f40/system-administrators-guide/package-management/DNF/) / [`Zypper`](https://de.opensuse.org/Zypper) package managers.

> Requires rpm version 4.15 or newer to support Ed25519 signatures.

> `$basearch` is supported by modern distributions. On older systems that don't expand `$basearch`, replace it in the `baseurl` with your architecture explicitly (for example, `.../rpm/cli/x86_64` or `.../rpm/cli/aarch64`).

##### Installation via DNF/YUM

1. Add the repository:

```shell
sudo tee /etc/yum.repos.d/stackit.repo > /dev/null << 'EOF'
[stackit]
name=STACKIT CLI
baseurl=https://packages.stackit.cloud/rpm/cli/$basearch
enabled=1
gpgcheck=1
gpgkey=https://packages.stackit.cloud/keys/key.gpg
EOF
```

2. Install the CLI:

```shell
sudo dnf install stackit
```

##### Installation via Zypper

1. Add the repository:

```shell
sudo tee /etc/zypp/repos.d/stackit.repo > /dev/null << 'EOF'
[stackit]
name=STACKIT CLI
baseurl=https://packages.stackit.cloud/rpm/cli/$basearch
enabled=1
gpgcheck=1
gpgkey=https://packages.stackit.cloud/keys/key.gpg
EOF
```

2. Install the CLI:

```shell
sudo zypper install stackit
```

#### Any distribution

Alternatively, you can install via [Homebrew](https://brew.sh/) or refer to one of the installation methods below.

> We are currently working on distributing the CLI on more package managers for Linux.

### Windows

> We are currently working on distributing the CLI on a package manager for Windows. For the moment, please refer to one of the installation methods below.

## Manual installation

You can also get the STACKIT CLI by compiling it from source or downloading a pre-compiled binary.

### Compile from source

1. Clone the repository
2. Build the application locally by running:

   ```bash
   make build
   ```

   To use the application from the root of the repository, you can run:

   ```bash
   ./bin/stackit <GROUP> <SUB-GROUP> <COMMAND> <ARGUMENT> <FLAGS>
   ```

3. Skip building and run the Go application directly using:

   ```bash
   go run . <GROUP> <SUB-GROUP> <COMMAND> <ARGUMENT> <FLAGS>
   ```

### FreeBSD

The STACKIT CLI can be installed through the [FreeBSD ports or packages](https://docs.freebsd.org/en/books/handbook/ports/).

To install the port:

```shell
cd /usr/ports/sysutils/stackit/ && make install clean
```

To add the package, run one of these commands:

```shell
pkg install sysutils/stackit
# OR
pkg install stackit
```

### Pre-compiled binary

1. Download the binary corresponding to your operating system and CPU architecture from our [Releases](https://github.com/stackitcloud/stackit-cli/releases) page
2. Extract the contents of the file to your file system and move it to your preferred location (make sure the directory is added to your `PATH`)
