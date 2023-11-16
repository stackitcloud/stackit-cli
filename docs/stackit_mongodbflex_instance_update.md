## stackit mongodbflex instance update

Update a MongoDB Flex instance

### Synopsis

Update a MongoDB Flex instance.

```
stackit mongodbflex instance update INSTANCE_ID [flags]
```

### Examples

```
  Update the name of a MongoDB Flex instance
  $ stackit mongodbflex instance update xxx --name my-new-name

  Update the version of a MongoDB Flex instance
  $ stackit mongodbflex instance update xxx --version 6.0
```

### Options

```
      --acl strings              List of IP networks in CIDR notation which are allowed to access this instance (default [])
      --backup-schedule string   Backup schedule
      --cpu int                  Number of CPUs
      --flavor-id string         ID of the flavor
  -h, --help                     Help for "stackit mongodbflex instance update"
  -n, --name string              Instance name
      --ram int                  Amount of RAM (in GB)
      --replicas int             Number of replicas
      --storage-class string     Storage class
      --storage-size int         Storage size (in GB)
      --type string              Instance type, one of ["Single" "Replica" "Sharded"]
      --version string           Version
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit mongodbflex instance](./stackit_mongodbflex_instance.md)	 - Provides functionality for MongoDB Flex instances

