## stackit mongodbflex user create

Creates a MongoDB Flex user

### Synopsis

Creates a MongoDB Flex user.
The password is only visible upon creation and cannot be retrieved later.
Alternatively, you can reset the password and access the new one by running:
  $ stackit mongodbflex user reset-password USER_ID --instance-id INSTANCE_ID

```
stackit mongodbflex user create [flags]
```

### Examples

```
  Create a MongoDB Flex user for instance with ID "xxx" and specify the username
  $ stackit mongodbflex user create --instance-id xxx --username johndoe --role read --database default

  Create a MongoDB Flex user for instance with ID "xxx" with an automatically generated username
  $ stackit mongodbflex user create --instance-id xxx --role read --database default
```

### Options

```
      --database string      The database inside the MongoDB instance that the user has access to. If it does not exist, it will be created once the user writes to it
  -h, --help                 Help for "stackit mongodbflex user create"
      --instance-id string   ID of the instance
      --role strings         Roles of the user, possible values are ["read" "readWrite"] (default [read])
      --username string      Username of the user. If not specified, a random username will be assigned
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

