## stackit secrets-manager user update

Updates the write privileges Secrets Manager user

### Synopsis

Updates the write privileges Secrets Manager user.

```
stackit secrets-manager user update USER_ID [flags]
```

### Examples

```
  Enable write access of a Secrets Manager user with ID "xxx" of instance with ID "yyy"
  $ stackit secrets-manager user update xxx --instance-id yyy --enable-write

  Disable write access of a Secrets Manager user with ID "xxx" of instance with ID "yyy"
  $ stackit secrets-manager user update xxx --instance-id yyy --disable-write
```

### Options

```
      --disable-write        Set the user to have read-only access.
      --enable-write         Set the user to have write access.
  -h, --help                 Help for "stackit secrets-manager user update"
      --instance-id string   ID of the instance
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit secrets-manager user](./stackit_secrets-manager_user.md)	 - Provides functionality for Secrets Manager users

