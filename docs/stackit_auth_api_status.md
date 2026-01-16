## stackit auth api status

Shows authentication status for the STACKIT Terraform Provider and SDK

### Synopsis

Shows authentication status for the STACKIT Terraform Provider and SDK, including whether you are authenticated and with which account.

```
stackit auth api status [flags]
```

### Examples

```
  Show authentication status for the STACKIT Terraform Provider and SDK
  $ stackit auth api status
```

### Options

```
  -h, --help   Help for "stackit auth api status"
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

