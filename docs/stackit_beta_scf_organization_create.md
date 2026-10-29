## stackit beta scf organization create

Creates a STACKIT Cloud Foundry organization

### Synopsis

Creates a STACKIT Cloud Foundry organization.

```
stackit beta scf organization create [flags]
```

### Examples

```
  Create a STACKIT Cloud Foundry organization with name "my-org"
  $ stackit beta scf organization create --name my-org

  Create a STACKIT Cloud Foundry organization with name "my-org" on platform "xxx"
  $ stackit beta scf organization create --name my-org --platform-id xxx
```

### Options

```
  -h, --help                 Help for "stackit beta scf organization create"
      --name string          Name of the organization
      --platform-id string   The UUID of the platform to create the organization on
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

