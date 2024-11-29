## stackit beta volume describe

Shows details of a volume

### Synopsis

Shows details of a volume.

```
stackit beta volume describe [flags]
```

### Examples

```
  Show details of a volume with ID "xxx"
  $ stackit beta volume describe xxx

  Show details of a volume with ID "xxx" in JSON format
  $ stackit beta volume describe xxx --output-format json
```

### Options

```
  -h, --help   Help for "stackit beta volume describe"
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

* [stackit beta volume](./stackit_beta_volume.md)	 - Provides functionality for volumes

