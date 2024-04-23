## stackit secrets-manager instance delete

Deletes a Secrets Manager instance

### Synopsis

Deletes a Secrets Manager instance.

```
stackit secrets-manager instance delete INSTANCE_ID [flags]
```

### Examples

```
  Delete a Secrets Manager instance with ID "xxx"
  $ stackit secrets-manager instance delete xxx
```

### Options

```
  -h, --help   Help for "stackit secrets-manager instance delete"
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

* [stackit secrets-manager instance](./stackit_secrets-manager_instance.md)	 - Provides functionality for Secrets Manager instances

