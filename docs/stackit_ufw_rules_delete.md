## stackit ufw rules delete

Deletes a UFW rule instance

### Synopsis

Deletes a STACKIT Unified Firewall (UFW) rule instance.

```
stackit ufw rules delete INSTANCE_ID [flags]
```

### Examples

```
  Delete a UFW rule instance with ID "xxx"
  $ stackit ufw rules delete xxx
```

### Options

```
  -h, --help   Help for "stackit ufw rules delete"
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

