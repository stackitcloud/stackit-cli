## stackit network-area region update

Updates a existing regional configuration for a STACKIT Network Area (SNA)

### Synopsis

Updates a existing regional configuration for a STACKIT Network Area (SNA).

```
stackit network-area region update [flags]
```

### Examples

```
  Update a regional configuration "eu02" for a STACKIT Network Area with ID "xxx" in organization with ID "yyy" with new ipv4-default-nameservers "8.8.8.8"
  $ stackit network-area region update --network-area-id xxx --region eu02 --organization-id yyy --ipv4-default-nameservers 8.8.8.8

  Update a regional configuration "eu02" for a STACKIT Network Area with ID "xxx" in organization with ID "yyy" with new ipv4-default-nameservers "8.8.8.8", using the set region config
  $ stackit config set --region eu02
  $ stackit network-area region update --network-area-id xxx --organization-id yyy --ipv4-default-nameservers 8.8.8.8

  Update a new regional configuration for a STACKIT Network Area with ID "xxx" in organization with ID "yyy", ipv4 network range "192.168.0.0/24", ipv4 transfer network "192.168.1.0/24", default prefix length "24", max prefix length "25" and min prefix length "20"
  $ stackit network-area region update --network-area-id xxx --organization-id yyy --ipv4-network-ranges 192.168.0.0/24 --ipv4-transfer-network 192.168.1.0/24 --region "eu02" --ipv4-default-prefix-length 24 --ipv4-max-prefix-length 25 --ipv4-min-prefix-length 20

  Update a new regional configuration for a STACKIT Network Area with ID "xxx" in organization with ID "yyy", ipv4 network range "192.168.0.0/24", ipv4 transfer network "192.168.1.0/24", default prefix length "24", max prefix length "25" and min prefix length "20"
  $ stackit network-area region update --network-area-id xxx --organization-id yyy --ipv4-network-ranges 192.168.0.0/24 --ipv4-transfer-network 192.168.1.0/24 --region "eu02" --ipv4-default-prefix-length 24 --ipv4-max-prefix-length 25 --ipv4-min-prefix-length 20
```

### Options

```
  -h, --help                               Help for "stackit network-area region update"
      --ipv4-default-nameservers strings   List of default DNS name server IPs
      --ipv4-default-prefix-length int     The default prefix length for networks in the network area
      --ipv4-max-prefix-length int         The maximum prefix length for networks in the network area
      --ipv4-min-prefix-length int         The minimum prefix length for networks in the network area
      --network-area-id string             STACKIT Network Area (SNA) ID
      --organization-id string             Organization ID
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

* [stackit network-area region](./stackit_network-area_region.md)	 - Provides functionality for regional configuration of STACKIT Network Area (SNA)

