## stackit beta network-area update

Updates a STACKIT Network Area (SNA)

### Synopsis

Updates a STACKIT Network Area (SNA) in an organization.

```
stackit beta network-area update AREA_ID [flags]
```

### Examples

```
  Update network area with ID "xxx" in organization with ID "yyy" with new name "network-area-1-new"
  $ stackit beta network-area update xxx --organization-id yyy --name network-area-1-new
```

### Options

```
      --default-prefix-length int   The default prefix length for networks in the network area
      --dns-name-servers strings    List of DNS name server IPs
  -h, --help                        Help for "stackit beta network-area update"
      --max-prefix-length int       The maximum prefix length for networks in the network area
      --min-prefix-length int       The minimum prefix length for networks in the network area
  -n, --name string                 Network area name
      --organization-id string      Organization ID
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

* [stackit beta network-area](./stackit_beta_network-area.md)	 - Provides functionality for STACKIT Network Area (SNA)

