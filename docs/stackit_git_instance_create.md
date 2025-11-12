## stackit git instance create

Creates STACKIT Git instance

### Synopsis

Create a STACKIT Git instance by name.

```
stackit git instance create [flags]
```

### Examples

```
  Create a instance with name 'my-new-instance'
  $ stackit git instance create --name my-new-instance

  Create a instance with name 'my-new-instance' and flavor
  $ stackit git instance create --name my-new-instance --flavor git-100

  Create a instance with name 'my-new-instance' and acl
  $ stackit git instance create --name my-new-instance --acl 1.1.1.1/1
```

### Options

```
      --acl strings     Acl for the instance.
      --flavor string   Flavor of the instance.
  -h, --help            Help for "stackit git instance create"
      --name string     The name of the instance.
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

* [stackit git instance](./stackit_git_instance.md)	 - Provides functionality for STACKIT Git instances

