## stackit sqlserverflex instance db-compatibility list

Lists supported SQLServer Flex database compatibilities of an instance

### Synopsis

Lists supported SQLServer Flex database compatibilities of an instance.

```
stackit sqlserverflex instance db-compatibility list [flags]
```

### Examples

```
  List SQLServer Flex database compatibilities of instance with ID "xxx"
  $ stackit sqlserverflex instance db-compatibility list --instance-id xxx

  List SQLServer Flex database compatibilities of instance with ID "xxx" in JSON format
  $ stackit sqlserverflex instance db-compatibility list --instance-id xxx --output-format json

  List up to 10 SQLServer Flex database compatibilities of instance with ID "xxx"
  $ stackit sqlserverflex instance db-compatibility list --instance-id xxx --limit 10
```

### Options

```
  -h, --help                 Help for "stackit sqlserverflex instance db-compatibility list"
      --instance-id string   Instance ID
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

* [stackit sqlserverflex instance db-compatibility](./stackit_sqlserverflex_instance_db-compatibility.md)	 - Provides functionality for SQLServer Flex database compatibilities

