## stackit config

Provides functionality for CLI configuration options

### Synopsis

Provides functionality for CLI configuration options.
You can set and unset different configuration options via the "stackit config set" and "stackit config unset" commands.

Additionally, you can configure the CLI to use different profiles, each with its own configuration.
Additional profiles can be configured via the "STACKIT_CLI_PROFILE" environment variable or using the "stackit config profile set PROFILE" and "stackit config profile unset" commands.
The environment variable takes precedence over what is set via the commands.

```
stackit config [flags]
```

### Options

```
  -h, --help   Help for "stackit config"
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit](./stackit.md)	 - Manage STACKIT resources using the command line
* [stackit config list](./stackit_config_list.md)	 - Lists the current CLI configuration values
* [stackit config profile](./stackit_config_profile.md)	 - Manage the CLI configuration profiles
* [stackit config set](./stackit_config_set.md)	 - Sets CLI configuration options
* [stackit config unset](./stackit_config_unset.md)	 - Unsets CLI configuration options

