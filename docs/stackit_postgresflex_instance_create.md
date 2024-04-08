## stackit postgresflex instance create

Creates a PostgreSQL Flex instance

### Synopsis

Creates a PostgreSQL Flex instance.

```
stackit postgresflex instance create [flags]
```

### Examples

```
  Create a PostgreSQL Flex instance with name "my-instance", ACL 0.0.0.0/0 (open access) and specify flavor by CPU and RAM. Other parameters are set to default values
  $ stackit postgresflex instance create --name my-instance --cpu 2 --ram 4 --acl 0.0.0.0/0

  Create a PostgreSQL Flex instance with name "my-instance", ACL 0.0.0.0/0 (open access) and specify flavor by ID. Other parameters are set to default values
  $ stackit postgresflex instance create --name my-instance --flavor-id xxx --acl 0.0.0.0/0

  Create a PostgreSQL Flex instance with name "my-instance", allow access to a specific range of IP addresses, specify flavor by CPU and RAM and set storage size to 20 GB. Other parameters are set to default values
  $ stackit postgresflex instance create --name my-instance --cpu 2 --ram 4 --acl 1.2.3.0/24 --storage-size 20
```

### Options

```
      --acl strings              The access control list (ACL). Must contain at least one valid subnet, for instance '0.0.0.0/0' for open access (discouraged), '1.2.3.0/24 for a public IP range of an organization, '1.2.3.4/32' for a single IP range, etc. (default [])
      --backup-schedule string   Backup schedule (default "0 0 * * *")
      --cpu int                  Number of CPUs
      --flavor-id string         ID of the flavor
  -h, --help                     Help for "stackit postgresflex instance create"
  -n, --name string              Instance name
      --ram int                  Amount of RAM (in GB)
      --storage-class string     Storage class (default "premium-perf2-stackit")
      --storage-size int         Storage size (in GB) (default 10)
      --type string              Instance type, one of ["Replica" "Single"] (default "Replica")
      --version string           PostgreSQL version. Defaults to the latest version available
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

* [stackit postgresflex instance](./stackit_postgresflex_instance.md)	 - Provides functionality for PostgreSQL Flex instances

