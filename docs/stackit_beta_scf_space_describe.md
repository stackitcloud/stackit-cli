## stackit beta scf space describe

Shows details of a STACKIT Cloud Foundry space

### Synopsis

Shows details of a STACKIT Cloud Foundry space.

```
stackit beta scf space describe SPACE_ID [flags]
```

### Examples

```
  Get details of the STACKIT Cloud Foundry space with ID "yyy" in the organization with ID "xxx"
  $ stackit beta scf space describe yyy --organization-id xxx

  Get details of the STACKIT Cloud Foundry space with ID "yyy" in the organization with ID "xxx" in JSON format
  $ stackit beta scf space describe yyy --organization-id xxx --output-format json
```

### Options

```
  -h, --help                     Help for "stackit beta scf space describe"
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

