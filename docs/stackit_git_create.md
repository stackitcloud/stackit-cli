## stackit git create

Creates STACKIT Git instance

### Synopsis

Create an STACKIT Git instance by name.

```
stackit git create [flags]
```

### Examples

```
  Create an instance with name 'my-new-instance'
  $ stackit git create --name my-new-instance
```

### Options

```
  -h, --help          Help for "stackit git create"
      --name string   The name of the instance.
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

* [stackit git](./stackit_git.md)	 - Provides functionality for STACKIT Git

