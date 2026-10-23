## stackit beta scf organization list

Lists all STACKIT Cloud Foundry organizations

### Synopsis

Lists all STACKIT Cloud Foundry organizations.

```
stackit beta scf organization list [flags]
```

### Examples

```
  List all STACKIT Cloud Foundry organizations
  $ stackit beta scf organization list

  List up to 10 STACKIT Cloud Foundry organizations
  $ stackit beta scf organization list --limit 10
```

### Options

```
  -h, --help        Help for "stackit beta scf organization list"
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

* [stackit beta scf organization](./stackit_beta_scf_organization.md)	 - Provides functionality for STACKIT Cloud Foundry organizations

