## stackit ufw rules describe

Shows details of a UFW rule instance

### Synopsis

Shows details of a STACKIT Unified Firewall (UFW) rule instance.

```
stackit ufw rules describe INSTANCE_ID [flags]
```

### Examples

```
  Get details of a UFW rule instance with ID "xxx"
  $ stackit ufw rule instance describe xxx

  Get details of a UFW rule instance with ID "xxx" in JSON format
  $ stackit ufw rule instance describe xxx --output-format json
```

### Options

```
  -h, --help   Help for "stackit ufw rules describe"
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

* [stackit ufw rules](./stackit_ufw_rules.md)	 - Provides functionality for UFW rules

