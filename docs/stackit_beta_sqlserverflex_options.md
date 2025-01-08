## stackit beta sqlserverflex options

Lists SQL Server Flex options

### Synopsis

Lists SQL Server Flex options (flavors, versions and storages for a given flavor)
Pass one or more flags to filter what categories are shown.

```
stackit beta sqlserverflex options [flags]
```

### Examples

```
  List SQL Server Flex flavors options
  $ stackit beta sqlserverflex options --flavors

  List SQL Server Flex available versions
  $ stackit beta sqlserverflex options --versions

  List SQL Server Flex storage options for a given flavor. The flavor ID can be retrieved by running "$ stackit beta sqlserverflex options --flavors"
  $ stackit beta sqlserverflex options --storages --flavor-id <FLAVOR_ID>

  List SQL Server Flex user roles and database compatibilities for a given instance. The IDs of existing instances can be obtained by running "$ stackit beta sqlserverflex instance list"
  $ stackit beta sqlserverflex options --user-roles --db-compatibilities --instance-id <INSTANCE_ID>
```

### Options

```
      --db-collations        Lists supported database collations for a given instance
      --db-compatibilities   Lists supported database compatibilities for a given instance
      --flavor-id string     The flavor ID to show storages for. Only relevant when "--storages" is passed
      --flavors              Lists supported flavors
  -h, --help                 Help for "stackit beta sqlserverflex options"
      --instance-id string   The instance ID to show user roles, database collations and database compatibilities for. Only relevant when "--user-roles", "--db-collations" or "--db-compatibilities" is passed
      --storages             Lists supported storages for a given flavor
      --user-roles           Lists supported user roles for a given instance
      --versions             Lists supported versions
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

* [stackit beta sqlserverflex](./stackit_beta_sqlserverflex.md)	 - Provides functionality for SQLServer Flex

