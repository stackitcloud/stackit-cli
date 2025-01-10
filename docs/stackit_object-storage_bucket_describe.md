## stackit object-storage bucket describe

Shows details of an Object Storage bucket

### Synopsis

Shows details of an Object Storage bucket.

```
stackit object-storage bucket describe BUCKET_NAME [flags]
```

### Examples

```
  Get details of an Object Storage bucket with name "my-bucket"
  $ stackit object-storage bucket describe my-bucket

  Get details of an Object Storage bucket with name "my-bucket" in JSON format
  $ stackit object-storage bucket describe my-bucket --output-format json
```

### Options

```
  -h, --help   Help for "stackit object-storage bucket describe"
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

* [stackit object-storage bucket](./stackit_object-storage_bucket.md)	 - Provides functionality for Object Storage buckets

