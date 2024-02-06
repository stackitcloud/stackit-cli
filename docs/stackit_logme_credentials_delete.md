## stackit logme credentials delete

Deletes credentials of a LogMe instance

### Synopsis

Deletes credentials of a LogMe instance.

```
stackit logme credentials delete CREDENTIALS_ID [flags]
```

### Examples

```
  Delete credentials with ID "xxx" of LogMe instance with ID "yyy"
  $ stackit logme credentials delete xxx --instance-id yyy
```

### Options

```
  -h, --help                 Help for "stackit logme credentials delete"
      --instance-id string   Instance ID
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit logme credentials](./stackit_logme_credentials.md)	 - Provides functionality for LogMe credentials

