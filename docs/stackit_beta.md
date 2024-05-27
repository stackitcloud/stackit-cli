## stackit beta

Contains Beta STACKIT CLI commands

### Synopsis

Contains Beta STACKIT CLI commands.

```
stackit beta [flags]
```

### Examples

```
  See the currently available Beta commands
  $ stackit beta --help

  Execute a Beta command
  $ stackit beta MY_COMMAND
```

### Options

```
  -h, --help   Help for "stackit beta"
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit](./stackit.md)	 - Manage STACKIT resources using the command line
* [stackit beta sqlserverflex](./stackit_beta_sqlserverflex.md)	 - Provides functionality for SQLServer Flex

