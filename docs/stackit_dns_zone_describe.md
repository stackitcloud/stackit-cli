## stackit dns zone describe

Shows details of a DNS zone

### Synopsis

Shows details of a DNS zone.

```
stackit dns zone describe ZONE_ID [flags]
```

### Examples

```
  Get details of a DNS zone with ID "xxx"
  $ stackit dns zone describe xxx

  Get details of a DNS zone with ID "xxx" in JSON format
  $ stackit dns zone describe xxx --output-format json
```

### Options

```
  -h, --help   Help for "stackit dns zone describe"
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

* [stackit dns zone](./stackit_dns_zone.md)	 - Provides functionality for DNS zones

