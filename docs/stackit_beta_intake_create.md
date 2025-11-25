## stackit beta intake create

Creates a new Intake

### Synopsis

Creates a new Intake.

```
stackit beta intake create [flags]
```

### Examples

```
  Create a new Intake with required parameters
  $ stackit beta intake create --display-name my-intake --runner-id xxx --catalog-uri "http://dremio.example.com" --catalog-warehouse "my-warehouse"

  Create a new Intake with a description, labels, and Dremio authentication
  $ stackit beta intake create --display-name my-intake --runner-id xxx --description "Production intake" --labels "env=prod,team=billing" --catalog-uri "http://dremio.example.com" --catalog-warehouse "my-warehouse" --catalog-auth-type "dremio" --dremio-token-endpoint "https://auth.dremio.cloud/oauth/token" --dremio-pat "MY_TOKEN"

  Create a new Intake with manual partitioning by a date field
  $ stackit beta intake create --display-name my-partitioned-intake --runner-id xxx --catalog-uri "http://dremio.example.com" --catalog-warehouse "my-warehouse" --catalog-partitioning "manual" --catalog-partition-by "day(__intake_ts)"
```

### Options

```
      --catalog-auth-type string       Authentication type for the catalog (e.g., 'none', 'dremio')
      --catalog-namespace string       The namespace to which data shall be written (default: 'intake')
      --catalog-partition-by strings   List of Iceberg partitioning expressions. Only used when --catalog-partitioning is 'manual'
      --catalog-partitioning string    The target table's partitioning. One of 'none', 'intake-time', 'manual'
      --catalog-table-name string      The table name to identify the table in Iceberg
      --catalog-uri string             The URI to the Iceberg catalog endpoint
      --catalog-warehouse string       The Iceberg warehouse to connect to
      --description string             Description
      --display-name string            Display name
      --dremio-pat string              Dremio personal access token. Required if auth-type is 'dremio'
      --dremio-token-endpoint string   Dremio OAuth 2.0 token endpoint URL. Required if auth-type is 'dremio'
  -h, --help                           Help for "stackit beta intake create"
      --labels stringToString          Labels in key=value format, separated by commas. Example: --labels "key1=value1,key2=value2" (default [])
      --runner-id string               The UUID of the Intake Runner to use
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

* [stackit beta intake](./stackit_beta_intake.md)	 - Provides functionality for intake

