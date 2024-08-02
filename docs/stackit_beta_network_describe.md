## stackit beta network describe

Shows details of a network

### Synopsis

Shows details of a network.

```
stackit beta network describe [flags]
```

### Examples

```
  Show details of a network with ID "xxx"
  $ stackit beta network describe xxx

  Show details of a network with ID "xxx" in JSON format
  $ stackit beta network describe xxx --output-format json
```

### Options

```
  -h, --help   Help for "stackit beta network describe"
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

* [stackit beta network](./stackit_beta_network.md)	 - Provides functionality for Network

