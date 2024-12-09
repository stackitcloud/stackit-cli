## stackit beta security-group create

Create security groups

### Synopsis

Create security groups.

```
stackit beta security-group create [flags]
```

### Examples

```
  create a named group
  $ stackit beta security-group create --name my-new-group

  create a named group with labels
  $ stackit beta security-group create --name my-new-group --labels label1=value1,label2=value2
```

### Options

```
      --description string   an optional description of the security group. Must be <= 127 chars
  -h, --help                 Help for "stackit beta security-group create"
      --labels strings       Labels are key-value string pairs which can be attached to a network-interface. E.g. '--labels key1=value1,key2=value2,...'
      --name string          the name of the security group. Must be <= 63 chars
      --stateful             create a stateful or a stateless security group
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

* [stackit beta security-group](./stackit_beta_security-group.md)	 - Manage security groups

