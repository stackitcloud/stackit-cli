## stackit config list

Lists the current CLI configuration values

### Synopsis

Lists the current CLI configuration values, based on the following sources (in order of precedence):
- Environment variable
  The environment variable is the name of the setting, with underscores ("_") instead of dashes ("-") and the "STACKIT" prefix.
  Example: you can set the project ID by setting the environment variable STACKIT_PROJECT_ID.
- Configuration set in CLI
  These are set using the "stackit config set" command
  Example: you can set the project ID by running "stackit config set --project-id xxx"

```
stackit config list [flags]
```

### Examples

```
  List your active configuration
  $ stackit config list
```

### Options

```
  -h, --help   Help for "stackit config list"
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit config](./stackit_config.md)	 - Provides functionality for CLI configuration options

