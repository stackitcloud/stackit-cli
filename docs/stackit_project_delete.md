## stackit project delete

Deletes a STACKIT project

### Synopsis

Deletes a STACKIT project.

```
stackit project delete [flags]
```

### Examples

```
  Delete the configured STACKIT project
  $ stackit project delete

  Delete a STACKIT project by explicitly providing the project ID
  $ stackit project delete --project-id xxx
```

### Options

```
  -h, --help   Help for "stackit project delete"
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

* [stackit project](./stackit_project.md)	 - Provides functionality regarding projects

