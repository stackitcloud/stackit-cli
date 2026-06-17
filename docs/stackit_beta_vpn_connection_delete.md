## stackit beta vpn connection delete

Deletes a VPN connection

### Synopsis

Deletes a VPN connection.

```
stackit beta vpn connection delete CONNECTION_ID [flags]
```

### Examples

```
  Delete a VPN connection
  $ stackit beta vpn connection delete xxx --gateway-id yyy
```

### Options

```
      --gateway-id string   Gateway ID
  -h, --help                Help for "stackit beta vpn connection delete"
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

