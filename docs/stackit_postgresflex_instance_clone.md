## stackit postgresflex instance clone

Clones a PostgreSQL Flex instance

### Synopsis

Clones a PostgreSQL Flex instance from a selected point in time.

```
stackit postgresflex instance clone INSTANCE_ID [flags]
```

### Examples

```
  Clone a PostgreSQL Flex instance with ID "xxx" . The recovery timestamp should be specified in UTC time following the format provided in the example.
  $ stackit postgresflex instance clone xxx --recovery-timestamp 2023-04-17T09:28:00+00:00

  Clone a PostgreSQL Flex instance with ID "xxx" from a selected recovery timestamp and specify storage class. If not specified, storage class from the existing instance will be used.
  $ stackit postgresflex instance clone xxx --recovery-timestamp 2023-04-17T09:28:00+00:00 --storage-class premium-perf6-stackit

  Clone a PostgreSQL Flex instance with ID "xxx" from a selected recovery timestamp and specify storage size. If not specified, storage size from the existing instance will be used.
  $ stackit postgresflex instance clone xxx --recovery-timestamp 2023-04-17T09:28:00+00:00 --storage-size 10
```

### Options

```
  -h, --help                        Help for "stackit postgresflex instance clone"
      --recovery-timestamp string   Recovery timestamp for the instance, in a date-time with the layout format, e.g. 2024-03-12T09:28:00+00:00
      --storage-class string        Storage class
      --storage-size int            Storage size (in GB)
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit postgresflex instance](./stackit_postgresflex_instance.md)	 - Provides functionality for PostgreSQL Flex instances

