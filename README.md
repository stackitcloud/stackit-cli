<div align="center">
<br>
<img src=".github/images/stackit-logo.png" alt="STACKIT logo" width="50%"/>
<br>
<br>
</div>

# STACKIT CLI (BETA)

[![Go Report Card](https://goreportcard.com/badge/github.com/stackitcloud/stackit-cli)](https://goreportcard.com/report/github.com/stackitcloud/stackit-cli) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/stackitcloud/stackit-cli) [![GitHub License](https://img.shields.io/github/license/stackitcloud/stackit-cli)](https://www.apache.org/licenses/LICENSE-2.0)

Welcome to the [STACKIT](https://www.stackit.de/en) CLI, a command-line interface for the STACKIT services.

This CLI is in a BETA state. More services and functionality will be supported soon.
Your feedback is appreciated!

<a name="warning-new-stackit-idp"></a>

> [!WARNING]
> On August 26 2024, The STACKIT Argus service was renamed to STACKIT Observability.
>
> This means that there is a new command group `observability`, which offers the same functionality as the deprecated `argus` command.
>
> Please make sure to **update your STACKIT CLI to the latest version after August 26 2024** to ensure that you start using `observability` command.

## Installation

Please refer to our [installation guide](./INSTALLATION.md) for instructions on how to install and get started using the STACKIT CLI.

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

Some commands are implemented at the root, group or subgroup level:

- `stackit config` to define variables to be used in future commands.
- `stackit ske enable` to enable the SKE engine on your project.

Help is available for any command by specifying the special flag `--help` (or simply `-h`):

- `stackit --help`
- `stackit -h`
- `stackit <GROUP> --help`
- `stackit <GROUP> <SUB-GROUP> --help`
- `stackit <GROUP> <SUB-GROUP> <COMMAND> --help`

## Available services

Below you can find a list of the STACKIT services already available in the CLI (along with their respective command names) and the ones that are currently planned to be integrated.

| Service                            | CLI Commands                                                                                                    | Status                    |
| ---------------------------------- | --------------------------------------------------------------------------------------------------------------- | ------------------------- |
| Observability                      | `observability`                                                                                                 | :white_check_mark:        |
| Infrastructure as a Service (IaaS) | `beta network-area` <br/> `beta network` <br/> `beta volume` <br/> `beta network-interface` <br/> `beta server` | :white_check_mark: (beta) |
| Authorization                      | `project`, `organization`                                                                                       | :white_check_mark:        |
| DNS                                | `dns`                                                                                                           | :white_check_mark:        |
| Kubernetes Engine (SKE)            | `ske`                                                                                                           | :white_check_mark:        |
| Load Balancer                      | `load-balancer`                                                                                                 | :white_check_mark:        |
| LogMe                              | `logme`                                                                                                         | :white_check_mark:        |
| MariaDB                            | `mariadb`                                                                                                       | :white_check_mark:        |
| MongoDB Flex                       | `mongodbflex`                                                                                                   | :white_check_mark:        |
| Object Storage                     | `object-storage`                                                                                                | :white_check_mark:        |
| OpenSearch                         | `opensearch`                                                                                                    | :white_check_mark:        |
| PostgreSQL Flex                    | `postgresflex`                                                                                                  | :white_check_mark:        |
| RabbitMQ                           | `rabbitmq`                                                                                                      | :white_check_mark:        |
| Redis                              | `redis`                                                                                                         | :white_check_mark:        |
| Resource Manager                   | `project`                                                                                                       | :white_check_mark:        |
| Secrets Manager                    | `secrets-manager`                                                                                               | :white_check_mark:        |
| Server Backup Management           | `beta server backup`                                                                                            | :white_check_mark: (beta) |
| Server Command (Run Command)       | `beta server command`                                                                                           | :white_check_mark: (beta) |
| Service Account                    | `service-account`                                                                                               | :white_check_mark:        |
| SQLServer Flex                     | `beta sqlserverflex`                                                                                            | :white_check_mark: (beta) |

## Authentication

Most of the commands will require you to be authenticated. Currently, it's possible to authenticate with your personal user or with a service account.

After successful authentication, the CLI stores credentials in your OS keychain. You won't need to log in again for the duration of your session, which is 2h by default but configurable by providing the `--session-time-limit` flag on the `config set` command (see [Configuration](#configuration)).

### Login with a personal user account

To authenticate as a user, run the command below and follow the steps in your browser.

```bash
stackit auth login
```

### Activate a service account

To authenticate using a service account, run:

```bash
stackit auth activate-service-account
```

For more details on how to set up authentication using a service account, check our [authentication guide](./AUTHENTICATION.md).

## Configuration

You can configure the CLI using the command:

```bash
stackit config
```

The configuration is saved in a file. The file's location varies depending on the operating system:

- Unix - `$XDG_CONFIG_HOME/stackit/cli-config.json`
- MacOS - `$HOME/Library/Application Support/stackit/cli-config.json`
- Windows - `%AppData%\stackit\cli-config.json`

The configuration options apply to all commands and can be set using the `stackit config set` command. For example, you can set a default `project-id` by running:

```bash
stackit config set --project-id xxxx-xxxx-xxxxx
```

To remove it, you can run:

```bash
stackit config unset --project-id
```

Run the `config set` command with the flag `--help` to get a list of all the available configuration options.

You can look up your current configuration by checking the configuration file or by running:

```bash
stackit config list
```

You can also edit the configuration file manually.

## Customization

### Pager

To specify a custom pager, use the `PAGER` environment variable.

If the variable is not set, STACKIT CLI uses the `less` as default pager.

When using `less` as a pager, STACKIT CLI will automatically pass following options

- -F, --quit-if-one-screen - Less will automatically exit if the entire file can be displayed on the first screen.
- -S, --chop-long-lines - Lines longer than the screen width will be chopped rather than being folded.
- -w, --hilite-unread - Temporarily highlights the first "new" line after a forward movement of a full page.
- -R, --RAW-CONTROL-CHARS - ANSI color and style sequences will be interpreted.

> These options will not be added automatically if a custom pager is defined.
>
> In that case, users can define the parameters by using the specific environment variable required by the `PAGER` (if supported).

> For example, if user sets the `PAGER` environment variable to `less` and would like to pass some arguments, `LESS` environment variable must be used as following:

> export PAGER="less"
>
> export LESS="-R"

## Autocompletion

If you wish to set up command autocompletion in your shell for the STACKIT CLI, please refer to our [autocompletion guide](./AUTOCOMPLETION.md).

## Reporting issues

If you encounter any issues or have suggestions for improvements, please open an issue in the [repository](https://github.com/stackitcloud/stackit-cli/issues).

## Contribute

Your contribution is welcome! For more details on how to contribute, refer to our [contribution guide](./CONTRIBUTION.md).

## License

Apache 2.0

## Useful Links

- [STACKIT Portal](https://portal.stackit.cloud/)

- [STACKIT](https://www.stackit.de/en/)

- [STACKIT Knowledge Base](https://docs.stackit.cloud/stackit/en/knowledge-base-85301704.html)

- [STACKIT Terraform Provider](https://registry.terraform.io/providers/stackitcloud/stackit/latest/docs)
