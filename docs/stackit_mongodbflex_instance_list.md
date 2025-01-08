## stackit mongodbflex instance list

Lists all MongoDB Flex instances

### Synopsis

Lists all MongoDB Flex instances.

```
stackit mongodbflex instance list [flags]
```

### Examples

```
  List all MongoDB Flex instances
  $ stackit mongodbflex instance list

  List all MongoDB Flex instances in JSON format
  $ stackit mongodbflex instance list --output-format json

  List up to 10 MongoDB Flex instances
  $ stackit mongodbflex instance list --limit 10
```

### Options

```
  -h, --help        Help for "stackit mongodbflex instance list"
      --limit int   Maximum number of entries to list
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

* [stackit mongodbflex instance](./stackit_mongodbflex_instance.md)	 - Provides functionality for MongoDB Flex instances

