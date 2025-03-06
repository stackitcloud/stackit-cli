## stackit image describe

Describes image

### Synopsis

Describes an image by its internal ID.

```
stackit image describe IMAGE_ID [flags]
```

### Examples

```
  Describe image "xxx"
  $ stackit image describe xxx
```

### Options

```
  -h, --help   Help for "stackit image describe"
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

