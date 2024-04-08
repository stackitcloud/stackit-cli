## stackit object-storage credentials delete

Deletes credentials of an Object Storage credentials group

### Synopsis

Deletes credentials of an Object Storage credentials group

```
stackit object-storage credentials delete CREDENTIALS_ID [flags]
```

### Examples

```
  Delete a credential with ID "xxx" of credentials group with ID "yyy"
  $ stackit object-storage credentials delete xxx --credentials-group-id yyy
```

### Options

```
      --credentials-group-id string   Credentials Group ID
  -h, --help                          Help for "stackit object-storage credentials delete"
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit object-storage credentials](./stackit_object-storage_credentials.md)	 - Provides functionality for Object Storage credentials

