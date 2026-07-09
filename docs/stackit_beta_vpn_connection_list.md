## stackit beta vpn connection list

Lists all VPN connections of a gateway

### Synopsis

Lists all VPN connections of a gateway.

```
stackit beta vpn connection list [flags]
```

### Examples

```
  List all VPN connections of a gateway
  $ stackit beta vpn connection list --gateway-id xxx
```

### Options

```
      --gateway-id string   Gateway ID
  -h, --help                Help for "stackit beta vpn connection list"
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

* [stackit beta vpn connection](./stackit_beta_vpn_connection.md)	 - Provides functionality for VPN connections

