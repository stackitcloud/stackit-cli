## stackit auth login

Logs in to the STACKIT CLI

### Synopsis

Logs in to the STACKIT CLI using a user account.
By default, the authentication uses a web-based authorization flow and opens a browser window where you can login to your STACKIT account. You can alternatively use the OAuth 2.0 device flow for environments where receiving a local callback is not possible.

```
stackit auth login [flags]
```

### Examples

```
  Login to the STACKIT CLI. This command will open a browser window where you can login to your STACKIT account
  $ stackit auth login

  Login to the STACKIT CLI using OAuth 2.0 device flow
  $ stackit auth login --use-device-flow
```

### Options

```
  -h, --help       Help for "stackit auth login"
      --port int   The port on which the callback server will listen to. By default, it tries to bind a port between 8000 and 8020.
                   When a value is specified, it will only try to use the specified port. Valid values are within the range of 8000 to 8020.
      --use-device-flow   Use OAuth 2.0 device authorization grant (device flow) instead of the browser callback flow.
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

