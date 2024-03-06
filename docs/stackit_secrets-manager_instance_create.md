## stackit secrets-manager instance create

Creates a Secrets Manager instance

### Synopsis

Creates a Secrets Manager instance.

```
stackit secrets-manager instance create [flags]
```

### Examples

```
  Create a Secrets Manager instance with name "my-instance"
  $ stackit secrets-manager instance create --name my-instance
```

### Options

```
  -h, --help          Help for "stackit secrets-manager instance create"
  -n, --name string   Instance name
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

