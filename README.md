# STACKIT CLI

Welcome to the STACKIT CLI, a command-line interface for the STACKIT services.

## Installation

To get started using it, you can:

1. Download the binary corresponding to your operating system and CPU architecture
2. Extract the contents of the file to your file system and move it to your preferred location (e.g. your home directory)
3. (For macOS only) Right click on the executable, select "Open". You will see a dialog stating the identity of the developer cannot be confirmed. Click on "Open" to allow the app to run on your Mac. We soon plan to certificate the STACKIT CLI to be trusted by macOS

Alternatively, you can use the STACKIT CLI by cloning the repository and either:

1. Build the application locally by running:

   ```bash
   $ go build -o ./bin/stackit
   ```

   To use the application from the root of the repository, you can run:

   ```bash
   $ ./bin/stackit <GROUP> <SUB-GROUP> <COMMAND> <ARGUMENT> <FLAGS>
   ```

2. Skip building and run the Go application directly using:

   ```bash
   $ go run . <GROUP> <SUB-GROUP> <COMMAND> <ARGUMENT> <FLAGS>
   ```

We also plan to integrate the STACKIT CLI on package managers such as APT and Brew.

## Usage

A typical command is structured as:

```
stackit <GROUP> <SUB-GROUP> <COMMAND> <ARGUMENT> <PARAMETER FLAGS> [OPTION FLAGS]
```

- `<GROUP>` can be the name of a service, such as `dns` or `mongodbflex`, or other groups for additional functionality, such as `config` to configure the CLI or `auth` to authenticate.
- `<SUB-GROUP>` should be the name (singular form) of a service resource, when `<GROUP>` is the name of a service. Examples: `zone`, `instance`.
- `<COMMAND>` is a command associated to the innermost group. Usually it's an action for the resource in question, such as `list` (to show all resources of the given type) or the CRUD operations `create`, `describe`, `update` and `delete`.
- `<ARGUMENT>` is required by some commands to specify a resource identifier. Examples: `stackit dns zone delete ZONE_ID`, `stackit ske cluster create CLUSTER_NAME`.
- `<PARAMETER FLAGS>` is a list of inputs necessary to execute the command, in the format `--[flag]` or `--[flag] [value]`. Some are required, while others are optional.
- `[OPTION FLAGS]` is a set of optional settings that modify the command's execution context. Examples: `--output-format=json` changes the format of the output to JSON, `--assume-yes` skips confirmation prompts.

Examples:

- `stackit ske cluster describe my-cluster --project-id xxx --output-format json`
- `stackit mongodbflex instance create --name my-instance --cpu 1 --ram 4 --acl 0.0.0.0/0 --assume-yes`
- `stackit dns zone delete my-zone`

Some commands are implemented at the root, group or sub-group level:

- `stackit config` to define variables to be used in future commands.
- `stackit ske enable` to enable the SKE engine on your project.

Help is available for any command by specifying the special flag `--help` (or simply `-h`):

- `stackit --help`
- `stackit -h`
- `stackit <GROUP> --help`
- `stackit <GROUP> <SUB-GROUP> --help`
- `stackit <GROUP> <SUB-GROUP> <COMMAND> --help`

## Authentication

Most of the commands will require you to be authenticated. Currently it's possible to authenticate with your personal user or with a service account.

After successful authentication, the CLI stores credentials in your OS keychain. You won't need to login again for the duration of your session, which is 2h by default but configurable by providing the `--session-time-limit` flag on the `config set` command (see [Configuration](#configuration)).

### Login with a personal user account

To authenticate as a user, run the command below and follow the steps in your browser.

```bash
$ stackit auth login
```

### Activate a service account

To authenticate using a service account, run:

```bash
$ stackit auth activate-service-account
```

For more details on how to setup authentication using a service account, check our [Authentication guide](./AUTHENTICATION.md)

## Configuration

You can configure the CLI using the command:

```bash
$ stackit config
```

The configurations are stored in `~/stackit/cli-config.json` and are valid for all commands. For example, you can set a default `project-id` by running:

```bash
$ stackit config set --project-id xxxx-xxxx-xxxxx
```

To remove it, you can run:

```bash
$ stackit config unset --project-id
```

Run the `config set` command with the flag `--help` to get a list of all of the available configuration options.

You can lookup your current configuration by checking the configuration file or by running:

```bash
$ stackit config list
```

You can also edit the configuration file manually.

## Reporting issues

If you encounter any issues or have suggestions for improvements, please reach out to the Developer Tools team or open a ticket through the [STACKIT Help Center](https://support.stackit.cloud/).

## Contribute (Coming soon to GitHub!)

Your contribution is welcome! For more details on how to contribute, refer to our [Contribution Guide](./CONTRIBUTION.md).

## License

Apache 2.0
