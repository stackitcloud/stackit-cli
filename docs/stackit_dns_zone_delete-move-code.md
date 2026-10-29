## stackit dns zone delete-move-code

Deletes the move code of a DNS zone

### Synopsis

Deletes (invalidates) the move code of a DNS zone, so it can no longer be used to move the zone to another project.

```
stackit dns zone delete-move-code ZONE_ID [flags]
```

### Examples

```
  Delete the move code of a DNS zone with ID "xxx"
  $ stackit dns zone delete-move-code xxx
```

### Options

```
  -h, --help   Help for "stackit dns zone delete-move-code"
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

