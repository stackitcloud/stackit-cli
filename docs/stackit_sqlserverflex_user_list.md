## stackit sqlserverflex user list

Lists all SQLServer Flex users of an instance

### Synopsis

Lists all SQLServer Flex users of an instance.

```
stackit sqlserverflex user list [flags]
```

### Examples

```
  List all SQLServer Flex users of instance with ID "xxx"
  $ stackit sqlserverflex user list --instance-id xxx

  List all SQLServer Flex users of instance with ID "xxx" in JSON format
  $ stackit sqlserverflex user list --instance-id xxx --output-format json

  List up to 10 SQLServer Flex users of instance with ID "xxx"
  $ stackit sqlserverflex user list --instance-id xxx --limit 10
```

### Options

```
  -h, --help                 Help for "stackit sqlserverflex user list"
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

* [stackit sqlserverflex user](./stackit_sqlserverflex_user.md)	 - Provides functionality for SQLServer Flex users

