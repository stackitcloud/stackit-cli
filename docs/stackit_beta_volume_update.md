## stackit beta volume update

Updates a volume

### Synopsis

Updates a volume.

```
stackit beta volume update VOLUME_ID [flags]
```

### Examples

```
  Update volume with ID "xxx" with new name "volume-1-new"
  $ stackit beta volume update xxx --name volume-1-new

  Update volume with ID "xxx" with new name "volume-1-new" and new description "volume-1-desc-new"
  $ stackit beta volume update xxx --name volume-1-new --description volume-1-desc-new

  Update volume with ID "xxx" with new name "volume-1-new", new description "volume-1-desc-new" and label(s)
  $ stackit beta volume update xxx --name volume-1-new --description volume-1-desc-new --labels key=value,foo=bar
```

### Options

```
      --description string      Volume description
  -h, --help                    Help for "stackit beta volume update"
      --labels stringToString   Labels are key-value string pairs which can be attached to a volume. E.g. '--labels key1=value1,key2=value2,...' (default [])
  -n, --name string             Volume name
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

* [stackit beta volume](./stackit_beta_volume.md)	 - Provides functionality for volumes

