## stackit object-storage credentials-group list

Lists all credentials groups that hold Object Storage access credentials

### Synopsis

Lists all credentials groups that hold Object Storage access credentials.

```
stackit object-storage credentials-group list [flags]
```

### Examples

```
  List all credentials groups
  $ stackit object-storage credentials-group list

  List all credentials groups in JSON format
  $ stackit object-storage credentials-group list --output-format json

  List up to 10 credentials groups
  $ stackit object-storage credentials-group list --limit 10
```

### Options

```
  -h, --help        Help for "stackit object-storage credentials-group list"
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

* [stackit object-storage credentials-group](./stackit_object-storage_credentials-group.md)	 - Provides functionality for Object Storage credentials group

