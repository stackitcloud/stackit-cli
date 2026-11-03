## stackit beta scf space delete

Deletes a STACKIT Cloud Foundry space

### Synopsis

Deletes a STACKIT Cloud Foundry space.

```
stackit beta scf space delete SPACE_ID [flags]
```

### Examples

```
  Delete the STACKIT Cloud Foundry space with ID "yyy" in the organization with ID "xxx"
  $ stackit beta scf space delete yyy --organization-id xxx
```

### Options

```
  -h, --help                     Help for "stackit beta scf space delete"
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

