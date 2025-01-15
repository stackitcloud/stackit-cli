## stackit beta image describe

Describes image

### Synopsis

Describes an image by its internal ID.

```
stackit beta image describe IMAGE_ID [flags]
```

### Examples

```
  Describe image "xxx"
  $ stackit beta image describe xxx
```

### Options

```
  -h, --help   Help for "stackit beta image describe"
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

* [stackit beta image](./stackit_beta_image.md)	 - Manage server images

