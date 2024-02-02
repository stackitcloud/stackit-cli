# Installation

## Package managers

### macOS

The STACKIT CLI can be installed through the [Homebrew](https://brew.sh/) package manager.

1. First, you need to register the [STACKIT tap](https://github.com/stackitcloud/homebrew-tap) via:

```shell
brew tap stackitcloud/tap
```

2. You can then install the CLI via:

```shell
brew install stackit
```

### Linux

We are working on distributing the CLI using a package manager for Linux. For the moment, you can either install via [Homebrew](https://brew.sh/) or refer to the [manual installation](#manual-installation) guide.

### Windows

We are working on distributing the CLI using a package manager for Windows. For the moment, please refer to the [manual installation](#manual-installation) guide.

## Manual installation

You can also get the STACKIT CLI by compiling it from source or downloading a pre-compiled binary.

### Compile from source

If you have Go 1.16+ installed, you can directly install via:

```shell
go install github.com/stackitcloud/stackit-cli@latest
```

> For more information, please refer to the [`go install` documentation](https://go.dev/ref/mod#go-install)

Alternativelly, you can:

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

### Pre-compiled binary

1. Download the binary corresponding to your operating system and CPU architecture from our [Releases](https://github.com/stackitcloud/stackit-cli/releases) page
2. Extract the contents of the file to your file system and move it to your preferred location (e.g. your home directory)
3. (For macOS only) Right click on the executable, select "Open". You will see a dialog stating the identity of the developer cannot be confirmed. Click on "Open" to allow the app to run on your Mac. We soon plan to certificate the STACKIT CLI to be trusted by macOS
