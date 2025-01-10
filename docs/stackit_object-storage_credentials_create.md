## stackit object-storage credentials create

Creates credentials for an Object Storage credentials group

### Synopsis

Creates credentials for an Object Storage credentials group. The credentials are only displayed upon creation, and will not be retrievable later.

```
stackit object-storage credentials create [flags]
```

### Examples

```
  Create credentials for a credentials group with ID "xxx"
  $ stackit object-storage credentials create --credentials-group-id xxx

  Create credentials for a credentials group with ID "xxx", including a specific expiration date
  $ stackit object-storage credentials create --credentials-group-id xxx --expire-date 2024-03-06T00:00:00.000Z
```

### Options

```
      --credentials-group-id string   Credentials Group ID
      --expire-date string            Expiration date for the credentials, in a date-time with the RFC3339 layout format, e.g. 2024-01-01T00:00:00Z
  -h, --help                          Help for "stackit object-storage credentials create"
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

* [stackit object-storage credentials](./stackit_object-storage_credentials.md)	 - Provides functionality for Object Storage credentials

