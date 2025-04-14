## stackit beta alb observability-credentials update

Update credentials

### Synopsis

Update credentials.

```
stackit beta alb observability-credentials update CREDENTIAL_REF_ARG [flags]
```

### Examples

```
  Update the password of observability credentials of Application Load Balancer with credentials reference "credentials-xxx", by providing the path to a file with the new password as flag
  $ stackit beta alb observability-credentials update credentials-xxx --username user1 --displayname user1 --password @./new-password.txt
```

### Options

```
  -d, --displayname string   Displayname for the credentials
  -h, --help                 Help for "stackit beta alb observability-credentials update"
      --password string      Password. Can be a string or a file path, if prefixed with "@" (example: @./password.txt).
  -u, --username string      Username for the credentials
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

