## stackit beta scf space create

Creates a STACKIT Cloud Foundry space

### Synopsis

Creates a STACKIT Cloud Foundry space in an organization.

```
stackit beta scf space create [flags]
```

### Examples

```
  Create a STACKIT Cloud Foundry space with name "my-space" in the organization with ID "xxx"
  $ stackit beta scf space create --organization-id xxx --name my-space
```

### Options

```
  -h, --help                     Help for "stackit beta scf space create"
      --name string              Name of the space
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

