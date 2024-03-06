## stackit secrets-manager user create

Creates a Secrets Manager user

### Synopsis

Creates a Secrets Manager user.
The username and password are auto-generated and provided upon creation.
A description can be provided to identify a user.

```
stackit secrets-manager user create [flags]
```

### Examples

```
  Create a Secrets Manager user for instance with ID "xxx" and description "yyy"
  $ stackit secrets-manager user create --instance-id xxx --description yyy

  Create a Secrets Manager user for instance with ID "xxx" and hides the generated password
  $ stackit secrets-manager user create --instance-id xxx --hide-password

  Create a Secrets Manager user for instance with ID "xxx" with write access to the secrets engine
  $ stackit secrets-manager user create --instance-id xxx --write
```

### Options

```
      --description string   A user chosen description to differentiate between multiple users
  -h, --help                 Help for "stackit secrets-manager user create"
      --hide-password        Hide password in output
      --instance-id string   ID of the instance
      --write                User write access to the secrets engine. If unset, user is read-only
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

