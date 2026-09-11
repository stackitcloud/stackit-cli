## stackit ufw rules list

Lists all UFW rules

### Synopsis

Lists all STACKIT Unified Firewall (UFW) rules.

```
stackit ufw rules list [flags]
```

### Examples

```
  List all UFW rules
  $ stackit ufw rules list

  List all UFW rules in JSON format
  $ stackit ufw rules list --output-format json

  List up to 10 UFW rules
  $ stackit ufw rules list --limit 10
```

### Options

```
  -h, --help        Help for "stackit ufw rules list"
      --limit int   Maximum number of entries to list
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

