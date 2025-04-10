## stackit beta alb credentials update

Update credentials

### Synopsis

Update credentials.

```
stackit beta alb credentials update CREDENTIAL_REF_ARG [flags]
```

### Examples

```
  Update the username
  $ stackit beta alb credentials update --username test-cred2 credentials-12345

  Update the displayname
  $ stackit beta alb credentials update --displayname new-name credentials-12345

  Update the password (is retrieved interactively or from ENV variable )
  $ stackit beta alb credentials update --password credentials-12345
```

### Options

```
  -d, --displayname string   the displayname for the credentials
  -h, --help                 Help for "stackit beta alb credentials update"
  -u, --username string      the username for the credentials
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

