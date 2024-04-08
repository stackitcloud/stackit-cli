## stackit mongodbflex instance delete

Deletes a MongoDB Flex instance

### Synopsis

Deletes a MongoDB Flex instance.

```
stackit mongodbflex instance delete INSTANCE_ID [flags]
```

### Examples

```
  Delete a MongoDB Flex instance with ID "xxx"
  $ stackit mongodbflex instance delete xxx
```

### Options

```
  -h, --help   Help for "stackit mongodbflex instance delete"
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

* [stackit mongodbflex instance](./stackit_mongodbflex_instance.md)	 - Provides functionality for MongoDB Flex instances

