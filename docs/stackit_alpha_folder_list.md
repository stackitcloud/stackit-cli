## stackit alpha folder list

Lists STACKIT folders

### Synopsis

Lists all STACKIT folders that match certain criteria.

```
stackit alpha folder list [flags]
```

### Examples

```
  List all STACKIT folders that are children of the configured project
  $ stackit folder list

  List all STACKIT folders that are children of a specific parent
  $ stackit folder list --parent-id xxx

  List all STACKIT folders that match the given project IDs, located under the same parent resource
  $ stackit folder list --project-id-like xxx,yyy,zzz

  List all STACKIT folders that a certain user is a member of
  $ stackit folder list --member example@email.com
```

### Options

```
      --creation-time-after string   Filter by creation timestamp, in a date-time with the RFC3339 layout format, e.g. 2023-01-01T00:00:00Z. The list of folders that were created after the given timestamp will be shown
  -h, --help                         Help for "stackit alpha folder list"
      --limit int                    Maximum number of entries to list
      --member string                Filter by member. The list of folders of which the member is part of will be shown
      --page-size int                Number of items fetched in each API call. Does not affect the number of items in the command output (default 50)
      --parent-id string             Filter by parent identifier
      --project-id-like strings      Filter by project identifier. Multiple project IDs can be provided, but they need to belong to the same parent resource (default [])
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, (one of: [json, pretty, none, yaml])
  -p, --project-id string      Project ID
      --region string          Target region for region-specific requests
      --verbosity string       Verbosity of the CLI, (one of: [debug, info, warning, error]) (default "info")
```

### SEE ALSO

* [stackit alpha folder](./stackit_alpha_folder.md)	 - Manages the lifecycle of folders.

