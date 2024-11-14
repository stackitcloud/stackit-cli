## stackit beta volume delete

Deletes a volume

### Synopsis

Deletes a volume.
If the volume is still in use, the deletion will fail


```
stackit beta volume delete [flags]
```

### Examples

```
  Delete volume with ID "xxx"
  $ stackit beta volume delete xxx
```

### Options

```
  -h, --help   Help for "stackit beta volume delete"
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

* [stackit beta volume](./stackit_beta_volume.md)	 - Provides functionality for Volume

