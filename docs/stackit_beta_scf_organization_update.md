## stackit beta scf organization update

Updates a STACKIT Cloud Foundry organization

### Synopsis

Updates a STACKIT Cloud Foundry organization. Only the specified fields are updated.

```
stackit beta scf organization update ORGANIZATION_ID [flags]
```

### Examples

```
  Update the name of a STACKIT Cloud Foundry organization with ID "xxx"
  $ stackit beta scf organization update xxx --name new-name

  Suspend a STACKIT Cloud Foundry organization with ID "xxx"
  $ stackit beta scf organization update xxx --suspended
```

### Options

```
  -h, --help          Help for "stackit beta scf organization update"
      --name string   Name of the organization
      --suspended     Whether the organization is suspended
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

