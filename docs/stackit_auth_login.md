## stackit auth login

Logs in to the STACKIT CLI

### Synopsis

Logs in to the STACKIT CLI using a user account.

```
stackit auth login [flags]
```

### Examples

```
  Login to the STACKIT CLI. This command will open a browser window where you can login to your STACKIT account
  $ stackit auth login
```

### Options

```
  -h, --help   Help for "stackit auth login"
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit auth](./stackit_auth.md)	 - Authenticates in the STACKIT CLI

