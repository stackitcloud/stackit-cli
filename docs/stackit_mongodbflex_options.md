## stackit mongodbflex options

Lists MongoDB Flex options

### Synopsis

Lists MongoDB Flex options (flavors, versions and storages for a given flavor)
Pass one or more flags to filter what categories are shown.

```
stackit mongodbflex options [flags]
```

### Examples

```
  List MongoDB Flex flavors options
  $ stackit mongodbflex options --flavors

  List MongoDB Flex available versions
  $ stackit mongodbflex options --versions

  List MongoDB Flex storage options for a given flavor. The flavor ID can be retrieved by running "$ stackit mongodbflex options --flavors"
  $ stackit mongodbflex options --storages --flavor-id <FLAVOR_ID>
```

### Options

```
      --flavor-id string   The flavor ID to show storages for. Only relevant when "--storages" is passed
      --flavors            Lists supported flavors
  -h, --help               Help for "stackit mongodbflex options"
      --storages           Lists supported storages for a given flavor
      --versions           Lists supported versions
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit mongodbflex](./stackit_mongodbflex.md)	 - Provides functionality for MongoDB Flex

