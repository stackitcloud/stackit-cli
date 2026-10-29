## stackit dns zone validate-move-code

Validates the move code of a DNS zone

### Synopsis

Validates the move code of a DNS zone. Checks that the provided code is valid, not expired, and belongs to the given zone.

```
stackit dns zone validate-move-code ZONE_ID [flags]
```

### Examples

```
  Validate the move code stored in "move-code.txt" of a DNS zone with ID "xxx"
  $ stackit dns zone validate-move-code xxx --code @./move-code.txt

  Validate a move code read from stdin for a DNS zone with ID "xxx"
  $ echo yyy | stackit dns zone validate-move-code xxx --code -
```

### Options

```
      --code secret   Code. Use @<file> to read from a file (example: @./secret.txt), or - to read from stdin (hidden input on a terminal). Literal values are not supported.
  -h, --help          Help for "stackit dns zone validate-move-code"
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

