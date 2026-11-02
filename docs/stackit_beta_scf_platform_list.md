## stackit beta scf platform list

Lists all STACKIT Cloud Foundry platforms

### Synopsis

Lists all STACKIT Cloud Foundry platforms. The number of fetched platforms is limited to 5000.

```
stackit beta scf platform list [flags]
```

### Examples

```
  List all STACKIT Cloud Foundry platforms
  $ stackit beta scf platform list

  List up to 10 STACKIT Cloud Foundry platforms
  $ stackit beta scf platform list --limit 10
```

### Options

```
  -h, --help        Help for "stackit beta scf platform list"
      --limit int   Maximum number of entries to list
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

* [stackit beta scf platform](./stackit_beta_scf_platform.md)	 - Provides functionality for STACKIT Cloud Foundry platforms

