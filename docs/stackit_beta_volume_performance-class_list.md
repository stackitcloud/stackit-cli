## stackit beta volume performance-class list

Lists all volume performance classes for a project

### Synopsis

Lists all volume performance classes for a project.

```
stackit beta volume performance-class list [flags]
```

### Examples

```
  Lists all volume performance classes
  $ stackit beta volume performance-class list

  Lists all volume performance classes which contains the label xxx
  $ stackit beta volume performance-class list --label-selector xxx

  Lists all volume performance classes in JSON format
  $ stackit beta volume performance-class list --output-format json

  Lists up to 10 volume performance classes
  $ stackit beta volume performance-class list --limit 10
```

### Options

```
  -h, --help                    Help for "stackit beta volume performance-class list"
      --label-selector string   Filter by label
      --limit int               Maximum number of entries to list
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
