## stackit beta sqlserverflex instance update

Updates a SQLServer Flex instance

### Synopsis

Updates a SQLServer Flex instance.

```
stackit beta sqlserverflex instance update INSTANCE_ID [flags]
```

### Examples

```
  Update the name of a SQLServer Flex instance with ID "xxx"
  $ stackit beta sqlserverflex instance update xxx --name my-new-name

  Update the backup schedule of a SQLServer Flex instance with ID "xxx"
  $ stackit beta sqlserverflex instance update xxx --backup-schedule "30 0 * * *"
```

### Options

```
      --acl strings              Lists of IP networks in CIDR notation which are allowed to access this instance (default [])
      --backup-schedule string   Backup schedule
      --cpu int                  Number of CPUs
      --flavor-id string         ID of the flavor
  -h, --help                     Help for "stackit beta sqlserverflex instance update"
  -n, --name string              Instance name
      --ram int                  Amount of RAM (in GB)
      --version string           Version
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

* [stackit beta sqlserverflex instance](./stackit_beta_sqlserverflex_instance.md)	 - Provides functionality for SQLServer Flex instances

