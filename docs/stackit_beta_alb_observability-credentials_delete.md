## stackit beta alb observability-credentials delete

Deletes credentials

### Synopsis

Deletes credentials.

```
stackit beta alb observability-credentials delete CREDENTIAL_REF [flags]
```

### Examples

```
  Delete credential with name "credential-12345"
  $ stackit beta alb observability-credentials delete credential-12345
```

### Options

```
  -h, --help   Help for "stackit beta alb observability-credentials delete"
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

* [stackit beta alb observability-credentials](./stackit_beta_alb_observability-credentials.md)	 - Provides functionality for application loadbalancer credentials

