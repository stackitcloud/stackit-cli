## stackit beta sqlserverflex instance create

Creates a SQLServer Flex instance

### Synopsis

Creates a SQLServer Flex instance.

```
stackit beta sqlserverflex instance create [flags]
```

### Examples

```
  Create a SQLServer Flex instance with name "my-instance" and specify flavor by ID.
  The flavor ID can be retrieved by running "$ stackit beta sqlserverflex flavor list"
  $ stackit beta sqlserverflex instance create --name my-instance --flavor-id xxx --backup-schedule "0 2 * * *" --retention-days 30 --storage-class premium-perf2-stackit --storage-size 10 --version 2022 --acl 1.2.3.0/24

  Create a SQLServer Flex instance with name "my-instance", specify flavor by ID, set storage size to 20 GB, and restrict access to a specific range of IP addresses.
  $ stackit beta sqlserverflex instance create --name my-instance --flavor-id xxx --storage-size 20 --backup-schedule "0 2 * * *" --retention-days 30 --storage-class premium-perf2-stackit --version 2022 --acl 1.2.3.0/24
```

### Options

```
      --acl strings                         The access control list (ACL). Must contain at least one valid subnet, for instance '0.0.0.0/0' for open access (discouraged), '1.2.3.0/24 for a public IP range of an organization, '1.2.3.4/32' for a single IP range, etc. (default [])
      --backup-schedule string              Backup schedule
      --edition string                      Edition of the SQLServer instance
      --encryption-kek-key-id string        The key identifier
      --encryption-kek-key-version string   The key version
      --encryption-kek-keyring-id string    The keyring identifier
      --encryption-service-account string   The service account
      --flavor-id string                    ID of the flavor. This flag will be required after 2027-02-28.
  -h, --help                                Help for "stackit beta sqlserverflex instance create"
  -n, --name string                         Instance name
      --retention-days int32                The days for how long the backup files should be stored before being cleaned up
      --storage-class string                Storage class
      --storage-size int                    Storage size (in GB)
      --version string                      SQLServer version
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

* [stackit beta sqlserverflex instance](./stackit_beta_sqlserverflex_instance.md)	 - Provides functionality for SQLServer Flex instances

