## stackit sqlserverflex database delete

Deletes a SQLServer Flex database

### Synopsis

Deletes a SQLServer Flex database.
This operation cannot be triggered asynchronously (the "--async" flag will have no effect).

```
stackit sqlserverflex database delete DATABASE_NAME [flags]
```

### Examples

```
  Delete a SQLServer Flex database with name "my-database" of instance with ID "xxx"
  $ stackit sqlserverflex database delete my-database --instance-id xxx
```

### Options

```
  -h, --help                 Help for "stackit sqlserverflex database delete"
      --instance-id string   SQLServer Flex instance ID
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

