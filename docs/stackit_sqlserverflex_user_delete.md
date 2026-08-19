## stackit sqlserverflex user delete

Deletes a SQLServer Flex user

### Synopsis

Deletes a SQLServer Flex user by ID. You can get the IDs of users for an instance by running:
  $ stackit sqlserverflex user list --instance-id <INSTANCE_ID>

```
stackit sqlserverflex user delete USER_ID [flags]
```

### Examples

```
  Delete a SQLServer Flex user with ID "xxx" for instance with ID "yyy"
  $ stackit sqlserverflex user delete xxx --instance-id yyy
```

### Options

```
  -h, --help                 Help for "stackit sqlserverflex user delete"
      --instance-id string   Instance ID
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, (one of: [json, pretty, none, yaml])
  -p, --project-id string      Project ID
      --region string          Target region for region-specific requests
      --verbosity string       Verbosity of the CLI, (one of: [debug, info, warning, error]) (default "info")
```

### SEE ALSO

* [stackit sqlserverflex user](./stackit_sqlserverflex_user.md)	 - Provides functionality for SQLServer Flex users

