## stackit security-group list

Lists security groups

### Synopsis

Lists security groups by its internal ID.

```
stackit security-group list [flags]
```

### Examples

```
  List all groups
  $ stackit security-group list

  List groups with labels
  $ stackit security-group list --label-selector label1=value1,label2=value2
```

### Options

```
  -h, --help                    Help for "stackit security-group list"
      --label-selector string   Filter by label
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

* [stackit security-group](./stackit_security-group.md)	 - Manage security groups

