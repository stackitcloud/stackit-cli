## stackit sqlserverflex instance user-role list

Lists supported SQLServer Flex user roles of an instance

### Synopsis

Lists supported SQLServer Flex user roles of an instance.

```
stackit sqlserverflex instance user-role list [flags]
```

### Examples

```
  List SQLServer Flex user roles of instance with ID "xxx"
  $ stackit sqlserverflex instance user-role list --instance-id xxx

  List SQLServer Flex user roles of instance with ID "xxx" in JSON format
  $ stackit sqlserverflex instance user-role list --instance-id xxx --output-format json

  List up to 10 SQLServer Flex user roles of instance with ID "xxx"
  $ stackit sqlserverflex instance user-role list --instance-id xxx --limit 10
```

### Options

```
  -h, --help                 Help for "stackit sqlserverflex instance user-role list"
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

* [stackit sqlserverflex instance user-role](./stackit_sqlserverflex_instance_user-role.md)	 - Provides functionality for SQLServer Flex user roles

