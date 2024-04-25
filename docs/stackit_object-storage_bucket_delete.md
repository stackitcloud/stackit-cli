## stackit object-storage bucket delete

Deletes an Object Storage bucket

### Synopsis

Deletes an Object Storage bucket.

```
stackit object-storage bucket delete BUCKET_NAME [flags]
```

### Examples

```
  Delete an Object Storage bucket with name "my-bucket"
  $ stackit object-storage bucket delete my-bucket
```

### Options

```
  -h, --help   Help for "stackit object-storage bucket delete"
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

