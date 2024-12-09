## stackit beta security-group delete

delete a security group

### Synopsis

delete a security group by its internal id

```
stackit beta security-group delete [flags]
```

### Examples

```
  delete a named group
  $ stackit beta security-group delete 43ad419a-c68b-4911-87cd-e05752ac1e31
```

### Options

```
  -h, --help   Help for "stackit beta security-group delete"
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

* [stackit beta security-group](./stackit_beta_security-group.md)	 - manage security groups.

