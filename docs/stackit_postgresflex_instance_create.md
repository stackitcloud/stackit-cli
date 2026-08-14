## stackit postgresflex instance create

Creates a PostgreSQL Flex instance

### Synopsis

Creates a PostgreSQL Flex instance.

```
stackit postgresflex instance create [flags]
```

### Examples

```
  Create a PostgreSQL Flex instance with name "my-instance", ACL 0.0.0.0/0 (open access).
  $ stackit postgresflex instance create --name my-instance --flavor-id xxx --acl 0.0.0.0/0 --storage-size 20 --retention-days 32 --version 17 --backup-schedule "6 6 * * *" --storage-size 10 --storage-class premium-perf2-stackit

  Create a PostgreSQL Flex instance with name "my-instance", allow access to a specific range of IP addresses.
  $ stackit postgresflex instance create --name my-instance --flavor-id xxx --acl 1.2.3.0/24 --storage-size 20 --retention-days 32 --version 17 --backup-schedule "6 6 * * *" --storage-size 10 --storage-class premium-perf2-stackit
```

### Options

```
      --access-scope string                 The access scope of the instance. It defines if the instance is public or airgapped. (one of: [PUBLIC, SNA])
      --acl strings                         The access control list (ACL). Must contain at least one valid subnet, for instance '0.0.0.0/0' for open access (discouraged), '1.2.3.0/24 for a public IP range of an organization, '1.2.3.4/32' for a single IP range, etc. (default [])
      --backup-schedule string              Backup schedule. This flag will be required after 2027-01-31.
      --encryption-kek-key-id string        The key identifier
      --encryption-kek-key-version string   The key version
      --encryption-kek-keyring-id string    The keyring identifier
      --encryption-service-account string   The service account
      --flavor-id string                    ID of the flavor. This flag will be required after 2027-01-31.
  -h, --help                                Help for "stackit postgresflex instance create"
  -n, --name string                         Instance name
      --retention-days int32                The days for how long the backup files should be stored before cleaned up (32 to 90). This flag will be required after 2027-01-31.
      --storage-class string                Storage class. This flag will be required after 2027-01-31.
      --storage-size int                    Storage size (in GB). This flag will be required after 2027-01-31.
      --version string                      PostgreSQL version. Defaults to the latest version available. This flag will be required after 2027-01-31.
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

* [stackit postgresflex instance](./stackit_postgresflex_instance.md)	 - Provides functionality for PostgreSQL Flex instances

