## stackit security-group describe

Describes security groups

### Synopsis

Describes security groups by its internal ID.

```
stackit security-group describe GROUP_ID [flags]
```

### Examples

```
  Describe group "xxx"
  $ stackit security-group describe xxx
```

### Options

```
  -h, --help   Help for "stackit security-group describe"
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --region string          Target region for region-specific requests
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit security-group](./stackit_security-group.md)	 - Manage security groups

