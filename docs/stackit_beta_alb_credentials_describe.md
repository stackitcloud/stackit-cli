## stackit beta alb credentials describe

Describes credentials

### Synopsis

Describes credentials.

```
stackit beta alb credentials describe CREDENTIAL_REF [flags]
```

### Examples

```
  Get details about credentials with name "credential-12345"
  $ stackit beta alb credential describe credential-12345
```

### Options

```
  -h, --help   Help for "stackit beta alb credentials describe"
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

* [stackit beta alb credentials](./stackit_beta_alb_credentials.md)	 - Provides functionality for application loadbalancer credentials

