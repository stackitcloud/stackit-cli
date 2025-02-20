## stackit volume performance-class describe

Shows details of a volume performance class

### Synopsis

Shows details of a volume performance class.

```
stackit volume performance-class describe VOLUME_PERFORMANCE_CLASS [flags]
```

### Examples

```
  Show details of a volume performance class with name "xxx"
  $ stackit volume performance-class describe xxx

  Show details of a volume performance class with name "xxx" in JSON format
  $ stackit volume performance-class describe xxx --output-format json
```

### Options

```
  -h, --help   Help for "stackit volume performance-class describe"
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

* [stackit volume performance-class](./stackit_volume_performance-class.md)	 - Provides functionality for volume performance classes available inside a project

