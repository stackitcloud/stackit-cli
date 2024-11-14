## stackit beta network-area route update

Updates a static route in a STACKIT Network Area (SNA)

### Synopsis

Updates a static route in a STACKIT Network Area (SNA).
This command is currently asynchonous only due to limitations in the waiting functionality of the SDK. This will be updated in a future release.


```
stackit beta network-area route update [flags]
```

### Examples

```
  Updates the label(s) of a static route with ID "xxx" in a STACKIT Network Area with ID "yyy" in organization with ID "zzz"
  $ stackit beta network-area route update xxx --labels key=value,foo=bar --organization-id yyy --network-area-id zzz
```

### Options

```
  -h, --help                     Help for "stackit beta network-area route update"
      --labels stringToString    Labels are key-value string pairs which can be attached to a route. A label can be provided with the format key=value and the flag can be used multiple times to provide a list of labels (default [])
      --network-area-id string   STACKIT Network Area ID
      --organization-id string   Organization ID
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

* [stackit beta network-area route](./stackit_beta_network-area_route.md)	 - Provides functionality for static routes in STACKIT Network Areas

