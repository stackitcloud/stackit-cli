## stackit secrets-manager user list

Lists all Secrets Manager users

### Synopsis

Lists all Secrets Manager users.

```
stackit secrets-manager user list [flags]
```

### Examples

```
  List all Secrets Manager users of instance with ID "xxx"
  $ stackit secrets-manager user list --instance-id xxx

  List all Secrets Manager users in JSON format with ID "xxx"
  $ stackit secrets-manager user list --instance-id xxx --output-format json

  List up to 10 Secrets Manager users with ID "xxx"
  $ stackit secrets-manager user list --instance-id xxx --limit 10
```

### Options

```
  -h, --help                 Help for "stackit secrets-manager user list"
      --instance-id string   Instance ID
      --limit int            Maximum number of entries to list
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit secrets-manager user](./stackit_secrets-manager_user.md)	 - Provides functionality for Secrets Manager users

