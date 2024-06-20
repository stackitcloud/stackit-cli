## stackit beta sqlserverflex instance create

Creates a SQLServer Flex instance

### Synopsis

Creates a SQLServer Flex instance.

```
stackit beta sqlserverflex instance create [flags]
```

### Examples

```
  Create a SQLServer Flex instance with name "my-instance" and specify flavor by CPU and RAM. Other parameters are set to default values
  $ stackit beta sqlserverflex instance create --name my-instance --cpu 1 --ram 4

  Create a SQLServer Flex instance with name "my-instance" and specify flavor by ID. Other parameters are set to default values
  $ stackit beta sqlserverflex instance create --name my-instance --flavor-id xxx

  Create a SQLServer Flex instance with name "my-instance", specify flavor by CPU and RAM, set storage size to 20 GB, and restrict access to a specific range of IP addresses. Other parameters are set to default values
  $ stackit beta sqlserverflex instance create --name my-instance --cpu 1 --ram 4 --storage-size 20  --acl 1.2.3.0/24
```

### Options

```
      --acl strings              The access control list (ACL). Must contain at least one valid subnet, for instance '0.0.0.0/0' for open access (discouraged), '1.2.3.0/24 for a public IP range of an organization, '1.2.3.4/32' for a single IP range, etc. (default [])
      --backup-schedule string   Backup schedule
      --cpu int                  Number of CPUs
      --edition string           Edition of the SQLServer instance
      --flavor-id string         ID of the flavor
  -h, --help                     Help for "stackit beta sqlserverflex instance create"
  -n, --name string              Instance name
      --ram int                  Amount of RAM (in GB)
      --retention-days int       The days for how long the backup files should be stored before being cleaned up
      --storage-class string     Storage class
      --storage-size int         Storage size (in GB)
      --version string           SQLServer version
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit beta sqlserverflex instance](./stackit_beta_sqlserverflex_instance.md)	 - Provides functionality for SQLServer Flex instances

