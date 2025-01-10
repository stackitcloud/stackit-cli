## stackit postgresflex instance clone

Clones a PostgreSQL Flex instance

### Synopsis

Clones a PostgreSQL Flex instance from a selected point in time. The new cloned instance will be an independent instance with the same settings as the original instance unless the flags are specified.

```
stackit postgresflex instance clone INSTANCE_ID [flags]
```

### Examples

```
  Clone a PostgreSQL Flex instance with ID "xxx" from a selected recovery timestamp.
  $ stackit postgresflex instance clone xxx --recovery-timestamp 2023-04-17T09:28:00+00:00

  Clone a PostgreSQL Flex instance with ID "xxx" from a selected recovery timestamp and specify storage class.
  $ stackit postgresflex instance clone xxx --recovery-timestamp 2023-04-17T09:28:00+00:00 --storage-class premium-perf6-stackit

  Clone a PostgreSQL Flex instance with ID "xxx" from a selected recovery timestamp and specify storage size.
  $ stackit postgresflex instance clone xxx --recovery-timestamp 2023-04-17T09:28:00+00:00 --storage-size 10
```

### Options

```
  -h, --help                        Help for "stackit postgresflex instance clone"
      --recovery-timestamp string   Recovery timestamp for the instance, in a date-time with the layout format YYYY-MM-DDTHH:mm:ss±HH:mm, e.g. 2006-01-02T15:04:05-07:00
      --storage-class string        Storage class. If not specified, storage class from the existing instance will be used.
      --storage-size int            Storage size (in GB). If not specified, storage size from the existing instance will be used.
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

* [stackit postgresflex instance](./stackit_postgresflex_instance.md)	 - Provides functionality for PostgreSQL Flex instances

