# Introduction

Welcome to the STACKIT CLI, a command-line interface for the STACKIT services.

# Getting Started

Currently, to use the STACKIT CLI you will have to clone the repository and either:

1. Build the application locally by running:

   ```bash
   $ go build -o ./bin/stackit
   ```

   To use the application from the root of the repository, you can run:

   ```bash
   $ ./bin/stackit [command] [subcommands] [flags]
   ```

2. Skip building and run the Go application directly using:

   ```bash
   $ go run . [command] [subcommands] [flags]
   ```

We will soon make this repository public and integrate a release pipeline that will build the executables for several operating systems on every new release. We also plan to integrate the STACKIT CLI on package managers such as APT and Brew.

# Authentication

Most of the commands will require you to be authenticated. Currently it's possile to authenticate using a user account or a service account.

To authenticate as a user, run the command below and follow the steps in your browser.

```bash
$ ./bin/stackit auth login
```

We will soon provide detailed instructions on how to authenticate as a service account.
