## stackit network-area create

Creates a STACKIT Network Area (SNA)

### Synopsis

Creates a STACKIT Network Area (SNA) in an organization.

```
stackit network-area create [flags]
```

### Examples

```
  Create a network area with name "network-area-1" in organization with ID "xxx"
  $ stackit network-area create --name network-area-1 --organization-id xxx"

  Create a network area with name "network-area-1" in organization with ID "xxx" with labels "key=value,key1=value1"
  $ stackit network-area create --name network-area-1 --organization-id xxx --labels key=value,key1=value1
```

### Options

```
  -h, --help                     Help for "stackit network-area create"
      --labels stringToString    Labels are key-value string pairs which can be attached to a network-area. E.g. '--labels key1=value1,key2=value2,...' (default [])
  -n, --name string              Network area name
      --organization-id string   Organization ID
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

