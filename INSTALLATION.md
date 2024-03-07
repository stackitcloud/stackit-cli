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

#### Debian/Ubuntu (`APT`)

The STACKIT CLI can be installed through the [`APT`](https://ubuntu.com/server/docs/package-management) package manager.

1. Import the STACKIT public key:

```shell
curl https://object.storage.eu01.onstackit.cloud/stackit-public-key/key.gpg | sudo gpg --dearmor -o /usr/share/keyrings/stackit.gpg
```

2. Add the STACKIT CLI package repository as a package source:

```shell
echo "deb [signed-by=/usr/share/keyrings/stackit.gpg] https://object.storage.eu01.onstackit.cloud/stackit-cli-apt stackit main" | sudo tee -a /etc/apt/sources.list.d/stackit.list
```

3. Update repository information and install the `stackit` package:

```shell
sudo apt-get update
sudo apt-get install stackit
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

### Pre-compiled binary

1. Download the binary corresponding to your operating system and CPU architecture from our [Releases](https://github.com/stackitcloud/stackit-cli/releases) page
2. Extract the contents of the file to your file system and move it to your preferred location (make sure the directory is added to your `PATH`)
