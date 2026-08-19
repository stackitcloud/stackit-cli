## stackit sqlserverflex database list

Lists all SQLServer Flex databases

### Synopsis

Lists all SQLServer Flex databases.

```
stackit sqlserverflex database list [flags]
```

### Examples

```
  List all SQLServer Flex databases of instance with ID "xxx"
  $ stackit sqlserverflex database list --instance-id xxx

  List all SQLServer Flex databases of instance with ID "xxx" in JSON format
  $ stackit sqlserverflex database list --instance-id xxx --output-format json

  List up to 10 SQLServer Flex databases of instance with ID "xxx"
  $ stackit sqlserverflex database list --instance-id xxx --limit 10
```

### Options

```
  -h, --help                 Help for "stackit sqlserverflex database list"
      --instance-id string   SQLServer Flex instance ID
      --limit int            Maximum number of entries to list
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

* [stackit sqlserverflex database](./stackit_sqlserverflex_database.md)	 - Provides functionality for SQLServer Flex databases

