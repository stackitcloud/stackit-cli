## stackit git instance list

Lists all instances of STACKIT Git.

### Synopsis

Lists all instances of STACKIT Git for the current project.

```
stackit git instance list [flags]
```

### Examples

```
  List all STACKIT Git instances
  $ stackit git instance list

  Lists up to 10 STACKIT Git instances
  $ stackit git instance list --limit=10
```

### Options

```
  -h, --help        Help for "stackit git instance list"
      --limit int   Limit the output to the first n elements
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

