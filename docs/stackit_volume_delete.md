## stackit volume delete

Deletes a volume

### Synopsis

Deletes a volume.
If the volume is still in use, the deletion will fail


```
stackit volume delete VOLUME_ID [flags]
```

### Examples

```
  Delete volume with ID "xxx"
  $ stackit volume delete xxx
```

### Options

```
  -h, --help   Help for "stackit volume delete"
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

* [stackit volume](./stackit_volume.md)	 - Provides functionality for volumes

