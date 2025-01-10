## stackit beta sqlserverflex database describe

Shows details of an SQLServer Flex database

### Synopsis

Shows details of an SQLServer Flex database.

```
stackit beta sqlserverflex database describe DATABASE_NAME [flags]
```

### Examples

```
  Get details of an SQLServer Flex database with name "my-database" of instance with ID "xxx"
  $ stackit beta sqlserverflex database describe my-database --instance-id xxx

  Get details of an SQLServer Flex database with name "my-database" of instance with ID "xxx" in JSON format
  $ stackit beta sqlserverflex database describe my-database --instance-id xxx --output-format json
```

### Options

```
  -h, --help                 Help for "stackit beta sqlserverflex database describe"
      --instance-id string   SQLServer Flex instance ID
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

* [stackit beta sqlserverflex database](./stackit_beta_sqlserverflex_database.md)	 - Provides functionality for SQLServer Flex databases

