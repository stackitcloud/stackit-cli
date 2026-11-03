## stackit beta scf organization quota list

Lists the quotas of a STACKIT Cloud Foundry organization

### Synopsis

Lists the quotas of a STACKIT Cloud Foundry organization.

```
stackit beta scf organization quota list ORGANIZATION_ID [flags]
```

### Examples

```
  List the quotas of a STACKIT Cloud Foundry organization with ID "xxx"
  $ stackit beta scf organization quota list xxx

  List the quotas of a STACKIT Cloud Foundry organization with ID "xxx" in JSON format
  $ stackit beta scf organization quota list xxx --output-format json
```

### Options

```
  -h, --help   Help for "stackit beta scf organization quota list"
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

