## stackit secrets-manager user delete

Deletes a Secrets Manager user

### Synopsis

Deletes a Secrets Manager user by ID. You can get the IDs of users for an instance by running:
  $ stackit secrets-manager user delete --instance-id <INSTANCE_ID>

```
stackit secrets-manager user delete USER_ID [flags]
```

### Examples

```
  Delete a Secrets Manager user with ID "xxx" for instance with ID "yyy"
  $ stackit secrets-manager user delete xxx --instance-id yyy
```

### Options

```
  -h, --help                 Help for "stackit secrets-manager user delete"
      --instance-id string   Instance ID
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

