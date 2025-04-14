## stackit beta alb observability-credentials add

Adds observability credentials to an application load balancer

### Synopsis

Adds observability credentials (username and password) to an application load balancer.  The credentials can be for Observability or another monitoring tool.

```
stackit beta alb observability-credentials add [flags]
```

### Examples

```
  Add observability credentials to a load balancer with username "xxx" and display name "yyy", providing the path to a file with the password as flag
  $ stackit beta alb observability-credentials add --username xxx --password @./password.txt --display-name yyy
```

### Options

```
  -d, --displayname string   Displayname for the credentials
  -h, --help                 Help for "stackit beta alb observability-credentials add"
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

