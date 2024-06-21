## stackit beta sqlserverflex database delete

Deletes a SQLServer Flex database

### Synopsis

Deletes a SQLServer Flex database.
This operation cannot be triggered asynchronously (the "--async" flag will have no effect).

```
stackit beta sqlserverflex database delete DATABASE_NAME [flags]
```

### Examples

```
  Delete a SQLServer Flex database with name "my-database" of instance with ID "xxx"
  $ stackit beta sqlserverflex database delete my-database --instance-id xxx
```

### Options

```
  -h, --help                 Help for "stackit beta sqlserverflex database delete"
      --instance-id string   SQLServer Flex instance ID
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

