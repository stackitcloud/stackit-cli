## stackit sqlserverflex instance delete

Deletes a SQLServer Flex instance

### Synopsis

Deletes a SQLServer Flex instance.

```
stackit sqlserverflex instance delete INSTANCE_ID [flags]
```

### Examples

```
  Delete a SQLServer Flex instance with ID "xxx"
  $ stackit sqlserverflex instance delete xxx
```

### Options

```
  -h, --help[=true|false]   Help for "stackit sqlserverflex instance delete"
```

### Options inherited from parent commands

```
  -y, --assume-yes[=true|false]   If set, skips all confirmation prompts
      --async[=true|false]        If set, runs the command asynchronously
  -o, --output-format string      Output format, (one of: [json, pretty, none, yaml])
  -p, --project-id string         Project ID
      --region string             Target region for region-specific requests
      --verbosity string          Verbosity of the CLI, (one of: [debug, info, warning, error]) (default "info")
```

### SEE ALSO

* [stackit sqlserverflex instance](./stackit_sqlserverflex_instance.md)	 - Provides functionality for SQLServer Flex instances

