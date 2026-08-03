## stackit beta vpn gateway list

Lists all vpn gateways

### Synopsis

Lists all vpn gateways.

```
stackit beta vpn gateway list [flags]
```

### Examples

```
  List all vpn gateways
  $ stackit beta vpn gateway list

  List up to 4 vpn gateways
  $ stackit beta vpn gateway list --limit 4
```

### Options

```
  -h, --help        Help for "stackit beta vpn gateway list"
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

* [stackit beta vpn gateway](./stackit_beta_vpn_gateway.md)	 - Provides functionality for VPN gateway

