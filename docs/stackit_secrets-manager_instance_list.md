## stackit secrets-manager instance list

Lists all Secrets Manager instances

### Synopsis

Lists all Secrets Manager instances.

```
stackit secrets-manager instance list [flags]
```

### Examples

```
  List all Secrets Manager instances
  $ stackit secrets-manager instance list

  List all Secrets Manager instances in JSON format
  $ stackit secrets-manager instance list --output-format json

  List up to 10 Secrets Manager instances
  $ stackit secrets-manager instance list --limit 10
```

### Options

```
  -h, --help        Help for "stackit secrets-manager instance list"
      --limit int   Maximum number of entries to list
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

* [stackit secrets-manager instance](./stackit_secrets-manager_instance.md)	 - Provides functionality for Secrets Manager instances

