## stackit secrets-manager instance describe

Shows details of a Secrets Manager instance

### Synopsis

Shows details of a Secrets Manager instance.

```
stackit secrets-manager instance describe INSTANCE_ID [flags]
```

### Examples

```
  Get details of a Secrets Manager instance with ID "xxx"
  $ stackit secrets-manager instance describe xxx

  Get details of a Secrets Manager instance with ID "xxx" in a table format
  $ stackit secrets-manager instance describe xxx --output-format pretty
```

### Options

```
  -h, --help   Help for "stackit secrets-manager instance describe"
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit secrets-manager instance](./stackit_secrets-manager_instance.md)	 - Provides functionality for Secrets Manager instances

