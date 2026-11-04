## stackit git instance update

Updates STACKIT Git instance

### Synopsis

Updates a STACKIT Git instance.

```
stackit git instance update INSTANCE_ID [flags]
```

### Examples

```
  Update the ACL of the instance with ID "xxx"
  $ stackit git instance update xxx --acl 1.2.3.0/24

  Update the ACL of the instance with ID "xxx" with multiple ranges
  $ stackit git instance update xxx --acl 1.2.3.0/24,4.5.6.0/24
```

### Options

```
      --acl strings   List of IP networks in CIDR notation which are allowed to access this instance. Replaces the existing ACL. (default [])
  -h, --help          Help for "stackit git instance update"
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, (one of: [json, pretty, none, yaml])
  -p, --project-id string      Project ID
      --region string          Target region for region-specific requests
      --verbosity string       Verbosity of the CLI, (one of: [debug, info, warning, error]) (default "info")
```

### SEE ALSO

* [stackit git instance](./stackit_git_instance.md)	 - Provides functionality for STACKIT Git instances

