## stackit auth api login

Logs in for the STACKIT Terraform Provider and SDK

### Synopsis

Logs in for the STACKIT Terraform Provider and SDK using a user account.
The authentication is done via a web-based authorization flow, where the command will open a browser window in which you can login to your STACKIT account.
The credentials are stored separately from the CLI authentication and will be used by the STACKIT Terraform Provider and SDK.

```
stackit auth api login [flags]
```

### Examples

```
  Login for the STACKIT Terraform Provider and SDK. This command will open a browser window where you can login to your STACKIT account
  $ stackit auth api login
```

### Options

```
  -h, --help   Help for "stackit auth api login"
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

* [stackit auth api](./stackit_auth_api.md)	 - Manages authentication for the STACKIT Terraform Provider and SDK

