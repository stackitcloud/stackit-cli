## stackit argus credentials delete

Deletes credentials of an Argus instance

### Synopsis

Deletes credentials of an Argus instance.

```
stackit argus credentials delete USERNAME [flags]
```

### Examples

```
  Delete credentials of username "xxx" for Argus instance with ID "yyy"
  $ stackit argus credentials delete xxx --instance-id yyy
```

### Options

```
  -h, --help                 Help for "stackit argus credentials delete"
      --instance-id string   Instance ID
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

* [stackit argus credentials](./stackit_argus_credentials.md)	 - Provides functionality for Argus credentials

