## stackit public-ip ranges list

Lists all STACKIT public-ip ranges

### Synopsis

Lists all STACKIT public-ip ranges.

```
stackit public-ip ranges list [flags]
```

### Examples

```
  Lists all STACKIT public-ip ranges
  $ stackit public-ip ranges list

  Lists all STACKIT public-ip ranges, piping to a tool like fzf for interactive selection
  $ stackit public-ip ranges list -o pretty | fzf

  Lists up to 10 STACKIT public-ip ranges
  $ stackit public-ip ranges list --limit 10
```

### Options

```
  -h, --help        Help for "stackit public-ip ranges list"
      --limit int   Maximum number of entries to list
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

* [stackit public-ip ranges](./stackit_public-ip_ranges.md)	 - Provides functionality for STACKIT public-ip ranges

