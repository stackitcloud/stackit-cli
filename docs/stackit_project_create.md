## stackit project create

Creates a STACKIT project

### Synopsis

Creates a STACKIT project.

```
stackit project create [flags]
```

### Examples

```
  Create a STACKIT project
  $ stackit project create --parent-id xxxx --name my-project

  Create a STACKIT project with a set of labels
  $ stackit project create --parent-id xxxx --name my-project --label key=value --label foo=bar
```

### Options

```
  -h, --help                   Help for "stackit project create"
      --label stringToString   Labels are key-value string pairs which can be attached to a project. A label can be provided with the format key=value and the flag can be used multiple times to provide a list of labels (default [])
      --name string            Project name
      --parent-id string       Parent resource identifier. Both container ID (user-friendly) and UUID are supported
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

* [stackit project](./stackit_project.md)	 - Provides functionality regarding projects

