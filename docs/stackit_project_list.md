## stackit project list

Lists STACKIT projects

### Synopsis

Lists all STACKIT projects that match certain criteria.

```
stackit project list [flags]
```

### Examples

```
  List all STACKIT projects that the authenticated user or service account is a member of
  $ stackit project list

  List all STACKIT projects that are children of a specific parent
  $ stackit project list --parent-id xxx

  List all STACKIT projects that match the given project IDs, located under the same parent resource
  $ stackit project list --project-id-like xxx,yyy,zzz

  List all STACKIT projects that a certain user is a member of
  $ stackit project list --member example@email.com
```

### Options

```
      --creation-time-after string   Filter by creation timestamp, in a date-time with the RFC3339 layout format, e.g. 2023-01-01T00:00:00Z. The list of projects that were created after the given timestamp will be shown
  -h, --help                         Help for "stackit project list"
      --limit int                    Maximum number of entries to list
      --member string                Filter by member. The list of projects of which the member is part of will be shown
      --page-size int                Number of items fetched in each API call. Does not affect the number of items in the command output (default 50)
      --parent-id string             Filter by parent identifier
      --project-id-like strings      Filter by project identifier. Multiple project IDs can be provided, but they need to belong to the same parent resource (default [])
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

* [stackit project](./stackit_project.md)	 - Provides functionality regarding projects

