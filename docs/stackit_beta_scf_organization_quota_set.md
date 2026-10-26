## stackit beta scf organization quota set

Applies a quota to a STACKIT Cloud Foundry organization

### Synopsis

Applies an existing quota to a STACKIT Cloud Foundry organization. The quota limits are not modified by this command.

```
stackit beta scf organization quota set ORGANIZATION_ID [flags]
```

### Examples

```
  Apply the quota with ID "yyy" to the STACKIT Cloud Foundry organization with ID "xxx"
  $ stackit beta scf organization quota set xxx --quota-id yyy
```

### Options

```
  -h, --help              Help for "stackit beta scf organization quota set"
      --quota-id string   The UUID of the quota to apply
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

* [stackit beta scf organization quota](./stackit_beta_scf_organization_quota.md)	 - Manages quotas of STACKIT Cloud Foundry organizations

