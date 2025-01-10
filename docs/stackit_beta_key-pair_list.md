## stackit beta key-pair list

Lists all key pairs

### Synopsis

Lists all key pairs.

```
stackit beta key-pair list [flags]
```

### Examples

```
  Lists all key pairs
  $ stackit beta key-pair list

  Lists all key pairs which contains the label xxx
  $ stackit beta key-pair list --label-selector xxx

  Lists all key pairs in JSON format
  $ stackit beta key-pair list --output-format json

  Lists up to 10 key pairs
  $ stackit beta key-pair list --limit 10
```

### Options

```
  -h, --help                    Help for "stackit beta key-pair list"
      --label-selector string   Filter by label
      --limit int               Number of key pairs to list
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

* [stackit beta key-pair](./stackit_beta_key-pair.md)	 - Provides functionality for SSH key pairs

