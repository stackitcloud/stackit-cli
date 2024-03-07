## stackit object-storage credentials list

Lists all credentials for an Object Storage credentials group

### Synopsis

Lists all credentials for an Object Storage credentials group.

```
stackit object-storage credentials list [flags]
```

### Examples

```
  List all credentials for a credentials group with ID "xxx"
  $ stackit object-storage credentials list --credentials-group-id xxx

  List all credentials for a credentials group with ID "xxx" in JSON format
  $ stackit object-storage credentials list --credentials-group-id xxx --output-format json

  List up to 10 credentials for a credentials group with ID "xxx"
  $ stackit object-storage credentials list --credentials-group-id xxx --limit 10
```

### Options

```
      --credentials-group-id string   Credentials Group ID
  -h, --help                          Help for "stackit object-storage credentials list"
      --limit int                     Maximum number of entries to list
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit object-storage credentials](./stackit_object-storage_credentials.md)	 - Provides functionality for Object Storage credentials

