## stackit beta security-group update

Updates a security group

### Synopsis

Updates a named security group

```
stackit beta security-group update GROUP_ID [flags]
```

### Examples

```
  Update the name of group "xxx"
  $ stackit beta security-group update xxx --name my-new-name

  Update the labels of group "xxx"
  $ stackit beta security-group update xxx --labels label1=value1,label2=value2
```

### Options

```
      --description string      An optional description of the security group.
  -h, --help                    Help for "stackit beta security-group update"
      --labels stringToString   Labels are key-value string pairs which can be attached to a network-interface. E.g. '--labels key1=value1,key2=value2,...' (default [])
      --name string             The name of the security group.
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

