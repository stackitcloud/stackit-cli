## stackit observability instance describe

Shows details of an Observability instance

### Synopsis

Shows details of an Observability instance.

```
stackit observability instance describe INSTANCE_ID [flags]
```

### Examples

```
  Get details of an Observability instance with ID "xxx"
  $ stackit observability instance describe xxx

  Get details of an Observability instance with ID "xxx" in JSON format
  $ stackit observability instance describe xxx --output-format json
```

### Options

```
  -h, --help   Help for "stackit observability instance describe"
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

* [stackit observability instance](./stackit_observability_instance.md)	 - Provides functionality for Observability instances

