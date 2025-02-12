## stackit auth get-access-token

Prints a short-lived access token.

### Synopsis

Prints a short-lived access token which can be used e.g. for API calls.

```
stackit auth get-access-token [flags]
```

### Examples

```
  Print a short-lived access token
  $ stackit auth get-access-token
```

### Options

```
  -h, --help   Help for "stackit auth get-access-token"
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

