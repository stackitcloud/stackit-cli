## stackit network-area region describe

Describes a regional configuration for a STACKIT Network Area (SNA)

### Synopsis

Describes a regional configuration for a STACKIT Network Area (SNA).

```
stackit network-area region describe [flags]
```

### Examples

```
  Describe a regional configuration "eu02" for a STACKIT Network Area with ID "xxx" in organization with ID "yyy"
  $ stackit network-area region describe --network-area-id xxx --region eu02 --organization-id yyy

  Describe a regional configuration "eu02" for a STACKIT Network Area with ID "xxx" in organization with ID "yyy", using the set region config
  $ stackit config set --region eu02
  $ stackit network-area region describe --network-area-id xxx --organization-id yyy
```

### Options

```
  -h, --help                     Help for "stackit network-area region describe"
      --network-area-id string   STACKIT Network Area (SNA) ID
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

* [stackit network-area region](./stackit_network-area_region.md)	 - Provides functionality for regional configuration of STACKIT Network Area (SNA)

