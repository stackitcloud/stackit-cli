## stackit beta security-group update

Update a security group

### Synopsis

Update a named security group

```
stackit beta security-group update [flags]
```

### Examples

```
  Update the name of a group
  $ stackit beta security-group update 541d122f-0a5f-4bb0-94b9-b1ccbd7ba776 --name my-new-name

  Update the labels of a group
  $ stackit beta security-group update 541d122f-0a5f-4bb0-94b9-b1ccbd7ba776 --labels label1=value1,label2=value2
```

### Options

```
      --description string   an optional description of the security group. Must be <= 127 chars
  -h, --help                 Help for "stackit beta security-group update"
      --labels strings       Labels are key-value string pairs which can be attached to a network-interface. E.g. '--labels key1=value1,key2=value2,...'
      --name string          the name of the security group. Must be <= 63 chars
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

