## stackit organization list

Lists all organizations

### Synopsis

Lists all organizations.

```
stackit organization list [flags]
```

### Examples

```
  Lists organizations for your user
  $ stackit organization list

  Lists the first 10 organizations
  $ stackit organization list --limit 10
```

### Options

```
  -h, --help        Help for "stackit organization list"
      --limit int   Maximum number of entries to list (default 50)
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

* [stackit organization](./stackit_organization.md)	 - Manages organizations

