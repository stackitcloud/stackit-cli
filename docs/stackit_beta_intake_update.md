## stackit beta intake update

Updates an Intake

### Synopsis

Updates an Intake. Only the specified fields are updated.

```
stackit beta intake update INTAKE_ID [flags]
```

### Examples

```
  Update the display name of an Intake with ID "xxx"
  $ stackit beta intake update xxx --runner-id yyy --display-name new-intake-name

  Update the catalog details for an Intake with ID "xxx"
  $ stackit beta intake update xxx --runner-id yyy --catalog-uri "http://new.uri" --catalog-warehouse "new-warehouse"
```

### Options

```
      --catalog-auth-type string       Authentication type for the catalog (e.g., 'none', 'dremio')
      --catalog-namespace string       The namespace to which data shall be written
      --catalog-table-name string      The table name to identify the table in Iceberg
      --catalog-uri string             The URI to the Iceberg catalog endpoint
      --catalog-warehouse string       The Iceberg warehouse to connect to
      --description string             Description
      --display-name string            Display name
      --dremio-pat string              Dremio personal access token
      --dremio-token-endpoint string   Dremio OAuth 2.0 token endpoint URL
  -h, --help                           Help for "stackit beta intake update"
      --labels stringToString          Labels in key=value format, separated by commas. Example: --labels "key1=value1,key2=value2". (default [])
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

