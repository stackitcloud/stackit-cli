## stackit beta sqlserverflex database create

Creates a SQLServer Flex database

### Synopsis

Creates a SQLServer Flex database.
This operation cannot be triggered asynchronously (the "--async" flag will have no effect).

```
stackit beta sqlserverflex database create DATABASE_NAME [flags]
```

### Examples

```
  Create a SQLServer Flex database with name "my-database" on instance with ID "xxx"
  $ stackit beta sqlserverflex database create my-database --instance-id xxx --owner some-username
```

### Options

```
  -h, --help                 Help for "stackit beta sqlserverflex database create"
      --instance-id string   SQLServer Flex instance ID
      --owner string         Username of the owner user
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit beta sqlserverflex database](./stackit_beta_sqlserverflex_database.md)	 - Provides functionality for SQLServer Flex databases

