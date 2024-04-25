## stackit object-storage bucket list

Lists all Object Storage buckets

### Synopsis

Lists all Object Storage buckets.

```
stackit object-storage bucket list [flags]
```

### Examples

```
  List all Object Storage buckets
  $ stackit object-storage bucket list

  List all Object Storage buckets in JSON format
  $ stackit object-storage bucket list --output-format json

  List up to 10 Object Storage buckets
  $ stackit object-storage bucket list --limit 10
```

### Options

```
  -h, --help        Help for "stackit object-storage bucket list"
      --limit int   Maximum number of entries to list
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit object-storage bucket](./stackit_object-storage_bucket.md)	 - Provides functionality for Object Storage buckets

