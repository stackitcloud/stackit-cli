## stackit beta volume resize

Resizes a volume

### Synopsis

Resizes a volume.

```
stackit beta volume resize VOLUME_ID [flags]
```

### Examples

```
  Resize volume with ID "xxx" with new size 10 GB
  $ stackit beta volume resize xxx --size 10
```

### Options

```
  -h, --help       Help for "stackit beta volume resize"
      --size int   Volume size (GB)
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

