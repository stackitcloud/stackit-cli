## stackit beta network-area delete

Deletes a network area

### Synopsis

Deletes a network area in an organization.

```
stackit beta network-area delete [flags]
```

### Examples

```
  Delete network area with ID "xxx" in organization with ID "yyy"
  $ stackit beta network-area delete xxx --organization-id yyy
```

### Options

```
  -h, --help                     Help for "stackit beta network-area delete"
      --organization-id string   Organization ID
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

* [stackit beta network-area](./stackit_beta_network-area.md)	 - Provides functionality for Network Area

