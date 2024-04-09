## stackit object-storage disable

Disables Object Storage for a project

### Synopsis

Disables Object Storage for a project. All buckets must be deleted beforehand.

```
stackit object-storage disable [flags]
```

### Examples

```
  Disable Object Storage functionality for your project.
  $ stackit object-storage disable
```

### Options

```
  -h, --help   Help for "stackit object-storage disable"
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit object-storage](./stackit_object-storage.md)	 - Provides functionality regarding Object Storage

