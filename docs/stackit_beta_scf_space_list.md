## stackit beta scf space list

Lists all STACKIT Cloud Foundry spaces of an organization

### Synopsis

Lists all STACKIT Cloud Foundry spaces of an organization. The number of fetched spaces is limited to 5000.

```
stackit beta scf space list [flags]
```

### Examples

```
  List all STACKIT Cloud Foundry spaces of the organization with ID "xxx"
  $ stackit beta scf space list --organization-id xxx

  List up to 10 STACKIT Cloud Foundry spaces of the organization with ID "xxx"
  $ stackit beta scf space list --organization-id xxx --limit 10
```

### Options

```
  -h, --help                     Help for "stackit beta scf space list"
      --limit int                Maximum number of entries to list
      --organization-id string   The UUID of the organization
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

* [stackit beta scf space](./stackit_beta_scf_space.md)	 - Provides functionality for STACKIT Cloud Foundry spaces

