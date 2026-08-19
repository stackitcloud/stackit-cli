## stackit sqlserverflex flavor list

Lists SQLServer Flex flavors

### Synopsis

Lists SQLServer Flex flavors.

```
stackit sqlserverflex flavor list [flags]
```

### Examples

```
  List SQLServer Flex flavors
  $ stackit sqlserverflex flavor list

  List up to 10 SQLServer Flex flavors
  $ stackit sqlserverflex flavor list --limit 10
```

### Options

```
  -h, --help        Help for "stackit sqlserverflex flavor list"
      --limit int   Maximum number of entries to list
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

* [stackit sqlserverflex flavor](./stackit_sqlserverflex_flavor.md)	 - Provides functionality for SQLServer Flex flavors

