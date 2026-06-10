## stackit server service-account attach

Attach a service account to a server

### Synopsis

Attach a service account to a server

```
stackit server service-account attach [flags]
```

### Examples

```
  Attach a service account with mail "xxx@sa.stackit.cloud" to a server with ID "yyy"
  $ stackit server service-account attach --service-account-email xxx@sa.stackit.cloud --server-id yyy
```

### Options

```
  -h, --help                           Help for "stackit server service-account attach"
  -s, --server-id string               Server ID
  -a, --service-account-email string   Service Account Email
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

* [stackit server service-account](./stackit_server_service-account.md)	 - Allows attaching/detaching service accounts to servers

