## stackit beta sqlserverflex instance delete

Deletes a SQLServer Flex instance

### Synopsis

Deletes a SQLServer Flex instance.

```
stackit beta sqlserverflex instance delete INSTANCE_ID [flags]
```

### Examples

```
  Delete a SQLServer Flex instance with ID "xxx"
  $ stackit beta sqlserverflex instance delete xxx
```

### Options

```
  -h, --help   Help for "stackit beta sqlserverflex instance delete"
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --region string          Target region for region-specific requests
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit beta sqlserverflex instance](./stackit_beta_sqlserverflex_instance.md)	 - Provides functionality for SQLServer Flex instances

