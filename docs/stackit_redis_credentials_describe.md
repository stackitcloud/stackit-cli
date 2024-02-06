## stackit redis credentials describe

Shows details of credentials of a Redis instance

### Synopsis

Shows details of credentials of a Redis instance. The password will be shown in plain text in the output.

```
stackit redis credentials describe CREDENTIALS_ID [flags]
```

### Examples

```
  Get details of credentials of a Redis instance with ID "xxx" from instance with ID "yyy"
  $ stackit redis credentials describe xxx --instance-id yyy

  Get details of credentials of a Redis instance with ID "xxx" from instance with ID "yyy" in a table format
  $ stackit redis credentials describe xxx --instance-id yyy --output-format pretty
```

### Options

```
  -h, --help                 Help for "stackit redis credentials describe"
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

* [stackit redis credentials](./stackit_redis_credentials.md)	 - Provides functionality for Redis credentials

