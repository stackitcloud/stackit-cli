## stackit config

Provides functionality for CLI configuration options

### Synopsis

Provides functionality for CLI configuration options
The configuration is stored in a file in the user's config directory, which is OS dependent.
Windows: %APPDATA%\stackit
Linux: $XDG_CONFIG_HOME/stackit or $HOME/.config/stackit
macOS: $HOME/Library/Application Support/stackit
The configuration file is named `cli-config.json` and is created automatically when setting a configuration option.

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
* [stackit config set](./stackit_config_set.md)	 - Sets CLI configuration options
* [stackit config unset](./stackit_config_unset.md)	 - Unsets CLI configuration options

