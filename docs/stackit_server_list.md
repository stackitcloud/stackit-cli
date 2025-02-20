## stackit server list

Lists all servers of a project

### Synopsis

Lists all servers of a project.

```
stackit server list [flags]
```

### Examples

```
  Lists all servers
  $ stackit server list

  Lists all servers which contains the label xxx
  $ stackit server list --label-selector xxx

  Lists all servers in JSON format
  $ stackit server list --output-format json

  Lists up to 10 servers
  $ stackit server list --limit 10
```

### Options

```
  -h, --help                    Help for "stackit server list"
      --label-selector string   Filter by label
      --limit int               Maximum number of entries to list
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

* [stackit server](./stackit_server.md)	 - Provides functionality for servers

