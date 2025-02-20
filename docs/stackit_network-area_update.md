## stackit network-area update

Updates a STACKIT Network Area (SNA)

### Synopsis

Updates a STACKIT Network Area (SNA) in an organization.

```
stackit network-area update AREA_ID [flags]
```

### Examples

```
  Update network area with ID "xxx" in organization with ID "yyy" with new name "network-area-1-new"
  $ stackit network-area update xxx --organization-id yyy --name network-area-1-new
```

### Options

```
      --default-prefix-length int   The default prefix length for networks in the network area
      --dns-name-servers strings    List of DNS name server IPs
  -h, --help                        Help for "stackit network-area update"
      --labels stringToString       Labels are key-value string pairs which can be attached to a network-area. E.g. '--labels key1=value1,key2=value2,...' (default [])
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

* [stackit network-area](./stackit_network-area.md)	 - Provides functionality for STACKIT Network Area (SNA)

