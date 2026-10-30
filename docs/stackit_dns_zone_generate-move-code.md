## stackit dns zone generate-move-code

Generates a move code for a DNS zone

### Synopsis

Generates a move code for a DNS zone, which is required to move the zone to another project.
The code is only shown once upon creation, so make sure to save it. By default, the code is valid for 30 days.

```
stackit dns zone generate-move-code ZONE_ID [flags]
```

### Examples

```
  Generate a move code for a DNS zone with ID "xxx"
  $ stackit dns zone generate-move-code xxx

  Generate a move code for a DNS zone with ID "xxx" that is valid for one hour
  $ stackit dns zone generate-move-code xxx --ttl 3600
```

### Options

```
  -h, --help        Help for "stackit dns zone generate-move-code"
      --ttl int32   Validity duration of the move code in seconds (min 1, max 2592000 / 30 days). If not set, the API default is used
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, (one of: [json, pretty, none, yaml])
  -p, --project-id string      Project ID
      --region string          Target region for region-specific requests
      --verbosity string       Verbosity of the CLI, (one of: [debug, info, warning, error]) (default "info")
```

### SEE ALSO

* [stackit dns zone](./stackit_dns_zone.md)	 - Provides functionality for DNS zones

