## stackit beta scf platform quota describe

Shows details of a STACKIT Cloud Foundry platform quota

### Synopsis

Shows details of a STACKIT Cloud Foundry platform quota.

```
stackit beta scf platform quota describe PLATFORM_ID [flags]
```

### Examples

```
  Get details of the quota with ID "yyy" of a STACKIT Cloud Foundry platform with ID "xxx"
  $ stackit beta scf platform quota describe xxx --quota-id yyy

  Get details of the quota with ID "yyy" of a STACKIT Cloud Foundry platform with ID "xxx" in JSON format
  $ stackit beta scf platform quota describe xxx --quota-id yyy --output-format json
```

### Options

```
  -h, --help              Help for "stackit beta scf platform quota describe"
      --quota-id string   The UUID of the quota
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

* [stackit beta scf platform quota](./stackit_beta_scf_platform_quota.md)	 - Shows quotas of STACKIT Cloud Foundry platforms

