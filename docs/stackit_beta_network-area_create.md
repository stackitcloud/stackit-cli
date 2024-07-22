## stackit beta network-area create

Creates a network area

### Synopsis

Creates a network area in an organization.

```
stackit beta network-area create [flags]
```

### Examples

```
  Create a network area with name "network-area-1" in organization with ID "org-1" with network ranges and a transfer network
  $ stackit beta network-area create --name network-area-1 --organization-id org-1 --network-ranges "1.1.1.0/24,192.123.1.0/24" --transfer-network "192.160.0.0/24"

  Create a network area with name "network-area-2" in organization with ID "org-2" with network ranges, transfer network and DNS name server
  $ stackit beta network-area create --name network-area-2 --organization-id org-2 --network-ranges "1.1.1.0/24,192.123.1.0/24" --transfer-network "192.160.0.0/24" --dns-name-servers "1.1.1.1"

  Create a network area with name "network-area-3" in organization with ID "org-3" with network ranges, transfer network and additional options
  $ stackit beta network-area create --name network-area-3 --organization-id org-3 --network-ranges "1.1.1.0/24,192.123.1.0/24" --transfer-network "192.160.0.0/24" --default-prefix-length 25 --max-prefix-length 29 --min-prefix-length 24
```

### Options

```
      --default-prefix-length int   The default prefix length for networks in the network area
      --dns-name-servers strings    List of DNS name server IPs
  -h, --help                        Help for "stackit beta network-area create"
      --max-prefix-length int       The maximum prefix length for networks in the network area
      --min-prefix-length int       The minimum prefix length for networks in the network area
  -n, --name string                 Network area name
      --network-ranges strings      List of network ranges (default [])
      --organization-id string      Organization ID
      --transfer-network string     Transfer network in CIDR notation
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit beta network-area](./stackit_beta_network-area.md)	 - Provides functionality for Network Area

