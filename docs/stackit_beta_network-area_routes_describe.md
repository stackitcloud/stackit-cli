## stackit beta network-area routes describe

Shows details of a static route in a STACKIT Network Area (SNA)

### Synopsis

Shows details of a static route in a STACKIT Network Area (SNA).

```
stackit beta network-area routes describe [flags]
```

### Examples

```
  Show details of a static route with ID "xxx" in a STACKIT Network Area with ID "yyy" in organization with ID "zzz"
  $ stackit beta network-area routes describe xxx --network-area-id yyy --organization-id zzz

  Show details of a static route with ID "xxx" in a STACKIT Network Area with ID "yyy" in organization with ID "zzz" in JSON format
  $ stackit beta network-area routes describe xxx --network-area-id yyy --organization-id zzz --output-format json
```

### Options

```
  -h, --help                     Help for "stackit beta network-area routes describe"
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

* [stackit beta network-area routes](./stackit_beta_network-area_routes.md)	 - Provides functionality for static routes in STACKIT Network Areas

