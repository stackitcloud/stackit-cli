## stackit object-storage credentials-group create

Creates a credentials group to hold Object Storage access credentials

### Synopsis

Creates a credentials group to hold Object Storage access credentials.

```
stackit object-storage credentials-group create [flags]
```

### Examples

```
  Create credentials group to hold Object Storage access credentials
  $ stackit object-storage credentials-group create --name example
```

### Options

```
  -h, --help          Help for "stackit object-storage credentials-group create"
      --name string   Name of the group holding credentials
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit object-storage credentials-group](./stackit_object-storage_credentials-group.md)	 - Provides functionality for Object Storage credentials group

