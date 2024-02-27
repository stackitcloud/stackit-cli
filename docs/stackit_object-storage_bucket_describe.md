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

  Get details of an Object Storage bucket with name "my-bucket" in a table format
  $ stackit object-storage bucket describe my-bucket --output-format pretty
```

### Options

```
  -h, --help   Help for "stackit object-storage bucket describe"
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit object-storage bucket](./stackit_object-storage_bucket.md)	 - Provides functionality for Object Storage buckets

