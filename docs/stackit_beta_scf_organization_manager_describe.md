## stackit beta scf organization manager describe

Shows details of the organization manager of a STACKIT Cloud Foundry organization

### Synopsis

Shows details of the organization manager of a STACKIT Cloud Foundry organization.

```
stackit beta scf organization manager describe ORGANIZATION_ID [flags]
```

### Examples

```
  Get details of the organization manager of a STACKIT Cloud Foundry organization with ID "xxx"
  $ stackit beta scf organization manager describe xxx

  Get details of the organization manager of a STACKIT Cloud Foundry organization with ID "xxx" in JSON format
  $ stackit beta scf organization manager describe xxx --output-format json
```

### Options

```
  -h, --help   Help for "stackit beta scf organization manager describe"
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

* [stackit beta scf organization manager](./stackit_beta_scf_organization_manager.md)	 - Manages the organization manager of STACKIT Cloud Foundry organizations

