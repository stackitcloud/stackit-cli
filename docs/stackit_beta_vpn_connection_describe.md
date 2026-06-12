## stackit beta vpn connection describe

Shows details of a VPN connection

### Synopsis

Shows details of a VPN connection.

```
stackit beta vpn connection describe CONNECTION_ID [flags]
```

### Examples

```
  Show details of a VPN connection
  $ stackit beta vpn connection describe xxx --gateway-id yyy
```

### Options

```
      --gateway-id string   Gateway ID
  -h, --help                Help for "stackit beta vpn connection describe"
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --region string          Target region for region-specific requests
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit beta vpn connection](./stackit_beta_vpn_connection.md)	 - Provides functionality for VPN connections

