# Installation

## Linux

We will soon distribute the STACKIT CLI via [Snap](https://snapcraft.io/). For the moment, please refer to the [manual installation](#manual-installation) guide.

## macOS

The STACKIT CLI is available to download and install through the [Homebrew](https://brew.sh/) package manager.

1. First, you need to register the [STACKIT tap](https://github.com/stackitcloud/homebrew-tap) via:

```shell
brew tap stackitcloud/tap
```

2. You can then install the CLI via:

```shell
brew install stackit-cli
```

## Windows

We will soon distribute the STACKIT CLI via [Chocolatey](https://chocolatey.org/). For the moment, please refer to the [manual installation](#manual-installation) guide.

## Manual installation

Alternatively, you can get the STACKIT CLI by downloading a pre-compiled binary or compiling it from source.

### Pre-compiled binary

1. Download the binary corresponding to your operating system and CPU architecture from our [Releases](https://github.com/stackitcloud/stackit-cli/releases) page
2. Extract the contents of the file to your file system and move it to your preferred location (e.g. your home directory)
3. (For macOS only) Right click on the executable, select "Open". You will see a dialog stating the identity of the developer cannot be confirmed. Click on "Open" to allow the app to run on your Mac. We soon plan to certificate the STACKIT CLI to be trusted by macOS

### Compile from source

1. Clone the repository
2. Build the application locally by running:

   ```bash
   $ make build
   ```

   To use the application from the root of the repository, you can run:

   ```bash
   $ ./bin/stackit <GROUP> <SUB-GROUP> <COMMAND> <ARGUMENT> <FLAGS>
   ```

3. Skip building and run the Go application directly using:

   ```bash
   $ go run . <GROUP> <SUB-GROUP> <COMMAND> <ARGUMENT> <FLAGS>
   ```
