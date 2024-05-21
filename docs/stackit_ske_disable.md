## stackit ske disable

Disables SKE for a project

### Synopsis

Disables STACKIT Kubernetes Engine (SKE) for a project. It will delete all associated clusters.

```
stackit ske disable [flags]
```

### Examples

```
  Disable SKE functionality for your project, deleting all associated clusters
  $ stackit ske disable
```

### Options

```
  -h, --help   Help for "stackit ske disable"
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

* [stackit ske](./stackit_ske.md)	 - Provides functionality for SKE

