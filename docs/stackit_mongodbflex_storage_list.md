## stackit mongodbflex storage list

Lists MongoDB Flex storages for a certain flavor

### Synopsis

Lists MongoDB Flex storages for a certain flavor.

```
stackit mongodbflex storage list [flags]
```

### Examples

```
  List MongoDB Flex storages for flavor with ID "xxx"
  $ stackit mongodbflex storage list --flavor-id xxx

  List MongoDB Flex storages for flavor with ID "xxx" in JSON format
  $ stackit mongodbflex storage list --flavor-id xxx --output-format json

  List up to 10 storages for flavor with ID "xxx"
  $ stackit mongodbflex storage list --flavor-id xxx --limit 10
```

### Options

```
      --flavor-id string   Flavor ID
  -h, --help               Help for "stackit mongodbflex storage list"
      --limit int          Maximum number of entries to list
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

* [stackit mongodbflex storage](./stackit_mongodbflex_storage.md)	 - Provides functionality for MongoDB Flex storages for a certain flavor

