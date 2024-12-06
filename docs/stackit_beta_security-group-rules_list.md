## stackit beta security-group-rules list

Lists all security group rules in a security group of a project

### Synopsis

Lists all security group rules in a security group of a project

```
stackit beta security-group-rules list [flags]
```

### Examples

```
  Lists all security group rules in security group with ID "xxx"
  $ stackit beta security-group-rule list --security-group-id xxx

  Lists all security group rules in security group with ID "xxx" in JSON format
  $ stackit beta security-group-rule list --security-group-id xxx --output-format json

  Lists up to 10 security group rules in security group with ID "xxx"
  $ stackit beta security-group-rule list --security-group-id xxx --limit 10
```

### Options

```
  -h, --help                       Help for "stackit beta security-group-rules list"
      --limit int                  Maximum number of entries to list
      --security-group-id string   The security group ID
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

* [stackit beta security-group-rules](./stackit_beta_security-group-rules.md)	 - Provides functionality for security group rules

