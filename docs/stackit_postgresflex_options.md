## stackit postgresflex options

Lists PostgreSQL Flex options

### Synopsis

Lists PostgreSQL Flex options (flavors, versions and storages for a given flavor)
Pass one or more flags to filter what categories are shown.

```
stackit postgresflex options [flags]
```

### Examples

```
  List PostgreSQL Flex flavors options
  $ stackit postgresflex options --flavors

  List PostgreSQL Flex available versions
  $ stackit postgresflex options --versions

  List PostgreSQL Flex storage options for a given flavor. The flavor ID can be retrieved by running "$ stackit postgresflex options --flavors"
  $ stackit postgresflex options --storages --flavor-id <FLAVOR_ID>
```

### Options

```
      --flavor-id string   The flavor ID to show storages for. Only relevant when "--storages" is passed
      --flavors            Lists supported flavors
  -h, --help               Help for "stackit postgresflex options"
      --storages           Lists supported storages for a given flavor
      --versions           Lists supported versions
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit postgresflex](./stackit_postgresflex.md)	 - Provides functionality for PostgreSQL Flex

