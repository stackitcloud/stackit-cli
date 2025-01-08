## stackit observability credentials delete

Deletes credentials of an Observability instance

### Synopsis

Deletes credentials of an Observability instance.

```
stackit observability credentials delete USERNAME [flags]
```

### Examples

```
  Delete credentials of username "xxx" for Observability instance with ID "yyy"
  $ stackit observability credentials delete xxx --instance-id yyy
```

### Options

```
  -h, --help                 Help for "stackit observability credentials delete"
      --instance-id string   Instance ID
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

* [stackit observability credentials](./stackit_observability_credentials.md)	 - Provides functionality for Observability credentials

