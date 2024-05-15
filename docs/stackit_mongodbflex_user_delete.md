## stackit mongodbflex user delete

Deletes a MongoDB Flex user

### Synopsis

Deletes a MongoDB Flex user by ID. You can get the IDs of users for an instance by running:
  $ stackit mongodbflex user list --instance-id <INSTANCE_ID>

```
stackit mongodbflex user delete USER_ID [flags]
```

### Examples

```
  Delete a MongoDB Flex user with ID "xxx" for instance with ID "yyy"
  $ stackit mongodbflex user delete xxx --instance-id yyy
```

### Options

```
  -h, --help                 Help for "stackit mongodbflex user delete"
      --instance-id string   Instance ID
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

* [stackit mongodbflex user](./stackit_mongodbflex_user.md)	 - Provides functionality for MongoDB Flex users

