## stackit beta scf platform quota list

Lists the quotas of a STACKIT Cloud Foundry platform

### Synopsis

Lists the quotas of a STACKIT Cloud Foundry platform.

```
stackit beta scf platform quota list PLATFORM_ID [flags]
```

### Examples

```
  List the quotas of a STACKIT Cloud Foundry platform with ID "xxx"
  $ stackit beta scf platform quota list xxx

  List the quotas of a STACKIT Cloud Foundry platform with ID "xxx" in JSON format
  $ stackit beta scf platform quota list xxx --output-format json
```

### Options

```
  -h, --help   Help for "stackit beta scf platform quota list"
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

