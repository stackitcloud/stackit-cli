## stackit beta sqlserverflex instance describe

Shows details  of a SQLServer Flex instance

### Synopsis

Shows details  of a SQLServer Flex instance.

```
stackit beta sqlserverflex instance describe INSTANCE_ID [flags]
```

### Examples

```
  Get details of a SQLServer Flex instance with ID "xxx"
  $ stackit beta sqlserverflex instance describe xxx

  Get details of a SQLServer Flex instance with ID "xxx" in JSON format
  $ stackit beta sqlserverflex instance describe xxx --output-format json
```

### Options

```
  -h, --help   Help for "stackit beta sqlserverflex instance describe"
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

* [stackit beta sqlserverflex instance](./stackit_beta_sqlserverflex_instance.md)	 - Provides functionality for SQLServer Flex instances

