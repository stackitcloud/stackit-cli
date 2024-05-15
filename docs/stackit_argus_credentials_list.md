## stackit argus credentials list

Lists the usernames of all credentials for an Argus instance

### Synopsis

Lists the usernames of all credentials for an Argus instance.

```
stackit argus credentials list [flags]
```

### Examples

```
  List the usernames of all credentials for an Argus instance with ID "xxx"
  $ stackit argus credentials list --instance-id xxx

  List the usernames of all credentials for an Argus instance in JSON format
  $ stackit argus credentials list --instance-id xxx --output-format json

  List the usernames of up to 10 credentials for an Argus instance
  $ stackit argus credentials list --instance-id xxx --limit 10
```

### Options

```
  -h, --help                 Help for "stackit argus credentials list"
      --instance-id string   Instance ID
      --limit int            Maximum number of entries to list
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

* [stackit argus credentials](./stackit_argus_credentials.md)	 - Provides functionality for Argus credentials

