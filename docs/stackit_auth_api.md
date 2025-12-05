## stackit auth api

Manages authentication for the STACKIT Terraform Provider and SDK

### Synopsis

Manages authentication for the STACKIT Terraform Provider and SDK.

These commands allow you to authenticate with your personal STACKIT account
and share the credentials with the STACKIT Terraform Provider and SDK.
This provides an alternative to using service accounts for local development.

```
stackit auth api [flags]
```

### Options

```
  -h, --help   Help for "stackit auth api"
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

* [stackit auth](./stackit_auth.md)	 - Authenticates the STACKIT CLI
* [stackit auth api get-access-token](./stackit_auth_api_get-access-token.md)	 - Prints a short-lived access token for the STACKIT Terraform Provider and SDK
* [stackit auth api login](./stackit_auth_api_login.md)	 - Logs in for the STACKIT Terraform Provider and SDK
* [stackit auth api logout](./stackit_auth_api_logout.md)	 - Logs out from the STACKIT Terraform Provider and SDK
* [stackit auth api status](./stackit_auth_api_status.md)	 - Shows authentication status for the STACKIT Terraform Provider and SDK

