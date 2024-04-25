## stackit project role list

Lists roles and permissions of a project

### Synopsis

Lists roles and permissions of a project.

```
stackit project role list [flags]
```

### Examples

```
  List all roles and permissions of a project
  $ stackit project role list --project-id xxx

  List all roles and permissions of a project in JSON format
  $ stackit project role list --project-id xxx --output-format json

  List up to 10 roles and permissions of a project
  $ stackit project role list --project-id xxx --limit 10
```

### Options

```
  -h, --help        Help for "stackit project role list"
      --limit int   Maximum number of entries to list
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit project role](./stackit_project_role.md)	 - Provides functionality regarding project roles

