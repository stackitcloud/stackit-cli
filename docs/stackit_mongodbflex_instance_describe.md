## stackit mongodbflex instance describe

Shows details  of a MongoDB Flex instance

### Synopsis

Shows details  of a MongoDB Flex instance.

```
stackit mongodbflex instance describe INSTANCE_ID [flags]
```

### Examples

```
  Get details of a MongoDB Flex instance with ID "xxx"
  $ stackit mongodbflex instance describe xxx

  Get details of a MongoDB Flex instance with ID "xxx" in a table format
  $ stackit mongodbflex instance describe xxx --output-format pretty
```

### Options

```
  -h, --help   Help for "stackit mongodbflex instance describe"
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit mongodbflex instance](./stackit_mongodbflex_instance.md)	 - Provides functionality for MongoDB Flex instances

