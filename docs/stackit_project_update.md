## stackit project update

Updates a STACKIT project

### Synopsis

Updates a STACKIT project.

```
stackit project update [flags]
```

### Examples

```
  Update the name of the configured STACKIT project
  $ stackit project update --name my-updated-project

  Add labels to the configured STACKIT project
  $ stackit project update --label key=value,foo=bar

  Update the name of a STACKIT project by explicitly providing the project ID
  $ stackit project update --name my-updated-project --project-id xxx
```

### Options

```
  -h, --help                   Help for "stackit project update"
      --label stringToString   Labels are key-value string pairs which can be attached to a project. A label can be provided with the format key=value and the flag can be used multiple times to provide a list of labels (default [])
      --name string            Project name
      --parent-id string       Parent resource identifier. Both container ID (user-friendly) and UUID are supported
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

