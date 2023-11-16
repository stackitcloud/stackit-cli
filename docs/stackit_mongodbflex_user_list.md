## stackit mongodbflex user list

List all MongoDB Flex users of an instance

### Synopsis

List all MongoDB Flex users of an instance.

```
stackit mongodbflex user list [flags]
```

### Examples

```
  List all MongoDB Flex users of instance with ID "xxx"
  $ stackit mongodbflex user list --instance-id xxx

  List all MongoDB Flex users of instance with ID "xxx" in JSON format
  $ stackit mongodbflex user list --instance-id xxx --output-format json

  List up to 10 MongoDB Flex users of instance with ID "xxx"
  $ stackit mongodbflex user list --instance-id xxx --limit 10
```

### Options

```
  -h, --help                 Help for "stackit mongodbflex user list"
      --instance-id string   Instance ID
      --limit int            Maximum number of entries to list
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

