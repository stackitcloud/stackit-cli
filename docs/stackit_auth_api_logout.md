## stackit auth api logout

Logs out from the STACKIT Terraform Provider and SDK

### Synopsis

Logs out from the STACKIT Terraform Provider and SDK. This does not affect CLI authentication.

```
stackit auth api logout [flags]
```

### Examples

```
  Log out from the STACKIT Terraform Provider and SDK
  $ stackit auth api logout
```

### Options

```
  -h, --help   Help for "stackit auth api logout"
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

