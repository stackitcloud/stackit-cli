## stackit image list

Lists images

### Synopsis

Lists images by their internal ID.

```
stackit image list [flags]
```

### Examples

```
  List all images
  $ stackit image list

  List images with label
  $ stackit image list --label-selector ARM64,dev

  List the first 10 images
  $ stackit image list --limit=10
```

### Options

```
  -h, --help                    Help for "stackit image list"
      --label-selector string   Filter by label
      --limit int               Limit the output to the first n elements
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

* [stackit image](./stackit_image.md)	 - Manage server images

