## stackit postgresflex user describe

Shows details of a PostgreSQL Flex user

### Synopsis

Shows details of a PostgreSQL Flex user.
The user password is hidden inside the "host" field and replaced with asterisks, as it is only visible upon creation. You can reset it by running:
  $ stackit postgresflex user reset-password USER_ID --instance-id INSTANCE_ID

```
stackit postgresflex user describe USER_ID [flags]
```

### Examples

```
  Get details of a PostgreSQL Flex user with ID "xxx" of instance with ID "yyy"
  $ stackit postgresflex user describe xxx --instance-id yyy

  Get details of a PostgreSQL Flex user with ID "xxx" of instance with ID "yyy" in JSON format
  $ stackit postgresflex user describe xxx --instance-id yyy --output-format json
```

### Options

```
  -h, --help                 Help for "stackit postgresflex user describe"
      --instance-id string   ID of the instance
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

* [stackit postgresflex user](./stackit_postgresflex_user.md)	 - Provides functionality for PostgreSQL Flex users

