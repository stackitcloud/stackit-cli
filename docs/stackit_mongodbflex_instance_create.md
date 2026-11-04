## stackit mongodbflex instance create

Creates a MongoDB Flex instance

### Synopsis

Creates a MongoDB Flex instance.

```
stackit mongodbflex instance create [flags]
```

### Examples

```
  Create a MongoDB Flex instance with name "my-instance", ACL 0.0.0.0/0 (open access).
  $ stackit mongodbflex instance create --name my-instance --flavor-id xxx --acl 0.0.0.0/0 --type Replica --storage-size 20 --version 8.0 --backup-schedule "6 6 * * *" --storage-size 10 --storage-class premium-perf2-mongodb

  Create a MongoDB Flex instance with name "my-instance", allow access to a specific range of IP addresses.
  $ stackit mongodbflex instance create --name my-instance --flavor-id xxx --acl 1.2.3.0/24 --type Replica --storage-size 20 --version 8.0 --backup-schedule "6 6 * * *" --storage-size 10 --storage-class premium-perf2-mongodb
```

### Options

```
      --acl strings              The access control list (ACL). Must contain at least one valid subnet, for instance '0.0.0.0/0' for open access (discouraged), '1.2.3.0/24 for a public IP range of an organization, '1.2.3.4/32' for a single IP range, etc. (default [])
      --backup-schedule string   Backup schedule. This flag will be required after 2027-03-07. (default "0 0/6 * * *")
      --flavor-id string         ID of the flavor. This flag will be required after 2027-03-07.
  -h, --help                     Help for "stackit mongodbflex instance create"
  -n, --name string              Instance name
      --storage-class string     Storage class. This flag will be required after 2027-03-07. (default "premium-perf2-mongodb")
      --storage-size int         Storage size (in GB). This flag will be required after 2027-03-07. (default 10)
      --type string              Instance type, (one of: [Replica, Sharded, Single]) (default "Replica")
      --version string           MongoDB version. Defaults to the latest version available. This flag will be required after 2027-03-07.
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, (one of: [json, pretty, none, yaml])
  -p, --project-id string      Project ID
      --region string          Target region for region-specific requests
      --verbosity string       Verbosity of the CLI, (one of: [debug, info, warning, error]) (default "info")
```

### SEE ALSO

* [stackit mongodbflex instance](./stackit_mongodbflex_instance.md)	 - Provides functionality for MongoDB Flex instances

