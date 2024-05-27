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
  $ stackit sqlserverflex options --flavors

  List SQL Server Flex available versions
  $ stackit sqlserverflex options --versions

  List SQL Server Flex storage options for a given flavor. The flavor ID can be retrieved by running "$ stackit sqlserverflex options --flavors"
  $ stackit sqlserverflex options --storages --flavor-id <FLAVOR_ID>
```

### Options

```
      --flavor-id string   The flavor ID to show storages for. Only relevant when "--storages" is passed
      --flavors            Lists supported flavors
  -h, --help               Help for "stackit beta sqlserverflex options"
      --storages           Lists supported storages for a given flavor
      --versions           Lists supported versions
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

* [stackit beta sqlserverflex](./stackit_beta_sqlserverflex.md)	 - Provides functionality for SQLServer Flex

