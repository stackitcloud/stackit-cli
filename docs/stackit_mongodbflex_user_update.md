## stackit mongodbflex user update

Updates a MongoDB Flex user

### Synopsis

Updates a MongoDB Flex user.

```
stackit mongodbflex user update USER_ID [flags]
```

### Examples

```
  Update the roles of a MongoDB Flex user with ID "xxx" of instance with ID "yyy"
  $ stackit mongodbflex user update xxx --instance-id yyy --role read
```

### Options

```
      --database string      The database inside the MongoDB instance that the user has access to. If it does not exist, it will be created once the user writes to it
  -h, --help                 Help for "stackit mongodbflex user update"
      --instance-id string   ID of the instance
      --role strings         Roles of the user, possible values are ["read" "readWrite"] (default [])
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit mongodbflex user](./stackit_mongodbflex_user.md)	 - Provides functionality for MongoDB Flex users

