## stackit beta volume update

Updates a volume

### Synopsis

Updates a volume.

```
stackit beta volume update [flags]
```

### Examples

```
  Update volume with ID "xxx" with new name "volume-1-new"
  $ stackit beta volume update xxx --name volume-1-new

  Update volume with ID "xxx" with new name "volume-1-new" and new description "volume-1-desc-new"
  $ stackit beta volume update xxx --name volume-1-new --description volume-1-desc-new
```

### Options

```
      --description string     Volume description
  -h, --help                   Help for "stackit beta volume update"
      --label stringToString   Labels are key-value string pairs which can be attached to a volume. A label can be provided with the format key=value and the flag can be used multiple times to provide a list of labels (default [])
  -n, --name string            Volume name
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

* [stackit beta volume](./stackit_beta_volume.md)	 - Provides functionality for Volume

