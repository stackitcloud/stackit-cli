## stackit beta network-area routes create

Creates a static route in a STACKIT Network Area (SNA)

### Synopsis

Creates a static route in a STACKIT Network Area (SNA).
This command is currently asynchonous only due to limitations in the waiting functionality of the SDK. This will be updated in a future release.


```
stackit beta network-area routes create [flags]
```

### Examples

```
  Create a static route with prefix "1.1.1.0/24" and next hop "1.1.1.1" in a STACKIT Network Area with ID "xxx" in organization with ID "yyy"
  $ stackit beta network-area routes create --organization-id yyy --network-area-id xxx --prefix 1.1.1.0/24 --next-hop 1.1.1.1
```

### Options

```
  -h, --help                     Help for "stackit beta network-area routes create"
      --network-area-id string   STACKIT Network Area ID
      --next-hop string          Next hop IP address. Must be a valid IPv4
      --organization-id string   Organization ID
      --prefix string            Static route prefix
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

* [stackit beta network-area routes](./stackit_beta_network-area_routes.md)	 - Provides functionality for static routes in STACKIT Network Areas

