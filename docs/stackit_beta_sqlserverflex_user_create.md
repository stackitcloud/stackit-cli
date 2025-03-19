## stackit beta sqlserverflex user create

Creates a SQLServer Flex user

### Synopsis

Creates a SQLServer Flex user for an instance.

The password is only visible upon creation and cannot be retrieved later.
Alternatively, you can reset the password and access the new one by running:
  $ stackit beta sqlserverflex user reset-password USER_ID --instance-id INSTANCE_ID
Please refer to https://docs.stackit.cloud/stackit/en/creating-logins-and-users-in-sqlserver-flex-instances-210862358.html for additional information.

The allowed user roles for your instance can be obtained by running:
  $ stackit beta sqlserverflex options --user-roles --instance-id INSTANCE_ID

```
stackit beta sqlserverflex user create [flags]
```

### Examples

```
  Create a SQLServer Flex user for instance with ID "xxx" and specify the username, role and database
  $ stackit beta sqlserverflex user create --instance-id xxx --username johndoe --roles "##STACKIT_DatabaseManager##"

  Create a SQLServer Flex user for instance with ID "xxx", specifying multiple roles
  $ stackit beta sqlserverflex user create --instance-id xxx --username johndoe --roles "##STACKIT_LoginManager##,##STACKIT_DatabaseManager##"
```

### Options

```
  -h, --help                 Help for "stackit beta sqlserverflex user create"
      --instance-id string   ID of the instance
      --roles strings        Roles of the user
      --username string      Username of the user
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

* [stackit beta sqlserverflex user](./stackit_beta_sqlserverflex_user.md)	 - Provides functionality for SQLServer Flex users

