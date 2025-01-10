## stackit secrets-manager user describe

Shows details of a Secrets Manager user

### Synopsis

Shows details of a Secrets Manager user.

```
stackit secrets-manager user describe USER_ID [flags]
```

### Examples

```
  Get details of a Secrets Manager user with ID "xxx" of instance with ID "yyy"
  $ stackit secrets-manager user describe xxx --instance-id yyy

  Get details of a Secrets Manager user with ID "xxx" of instance with ID "yyy" in JSON format
  $ stackit secrets-manager user describe xxx --instance-id yyy --output-format json
```

### Options

```
  -h, --help                 Help for "stackit secrets-manager user describe"
      --instance-id string   ID of the instance
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

* [stackit secrets-manager user](./stackit_secrets-manager_user.md)	 - Provides functionality for Secrets Manager users

