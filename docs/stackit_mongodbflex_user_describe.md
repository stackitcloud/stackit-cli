## stackit mongodbflex user describe

Shows details of a MongoDB Flex user

### Synopsis

Shows details of a MongoDB Flex user.
The user password is hidden inside the "host" field and replaced with asterisks, as it is only visible upon creation. You can reset it by running:
  $ stackit mongodbflex user reset-password USER_ID --instance-id INSTANCE_ID

```
stackit mongodbflex user describe USER_ID [flags]
```

### Examples

```
  Get details of a MongoDB Flex user with ID "xxx" of instance with ID "yyy"
  $ stackit mongodbflex user list xxx --instance-id yyy

  Get details of a MongoDB Flex user with ID "xxx" of instance with ID "yyy" in table format
  $ stackit mongodbflex user list xxx --instance-id yyy --output-format pretty
```

### Options

```
  -h, --help                 Help for "stackit mongodbflex user describe"
      --instance-id string   ID of the instance
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

* [stackit mongodbflex user](./stackit_mongodbflex_user.md)	 - Provides functionality for MongoDB Flex users

