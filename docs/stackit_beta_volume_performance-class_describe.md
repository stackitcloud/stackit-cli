## stackit beta volume performance-class describe

Shows details of a volume performance class

### Synopsis

Shows details of a volume performance class.

```
stackit beta volume performance-class describe [flags]
```

### Examples

```
  Show details of a volume performance class with name "xxx"
  $ stackit beta volume performance-class describe xxx

  Show details of a volume performance class with name "xxx" in JSON format
  $ stackit beta volume performance-class describe xxx --output-format json
```

### Options

```
  -h, --help   Help for "stackit beta volume performance-class describe"
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

* [stackit beta volume performance-class](./stackit_beta_volume_performance-class.md)	 - Provides functionality for volume performance classes available inside a project

