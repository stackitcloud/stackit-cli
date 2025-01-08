## stackit postgresflex instance update

Updates a PostgreSQL Flex instance

### Synopsis

Updates a PostgreSQL Flex instance.

```
stackit postgresflex instance update INSTANCE_ID [flags]
```

### Examples

```
  Update the name of a PostgreSQL Flex instance
  $ stackit postgresflex instance update xxx --name my-new-name

  Update the version of a PostgreSQL Flex instance
  $ stackit postgresflex instance update xxx --version 6.0
```

### Options

```
      --acl strings              List of IP networks in CIDR notation which are allowed to access this instance (default [])
      --backup-schedule string   Backup schedule
      --cpu int                  Number of CPUs
      --flavor-id string         ID of the flavor
  -h, --help                     Help for "stackit postgresflex instance update"
  -n, --name string              Instance name
      --ram int                  Amount of RAM (in GB)
      --storage-class string     Storage class
      --storage-size int         Storage size (in GB)
      --type string              Instance type, one of ["Replica" "Single"]
      --version string           Version
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

