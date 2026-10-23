## stackit dns zone move

Moves a DNS zone to another project

### Synopsis

Moves a DNS zone to another project.
The API call this command makes is always performed against the destination project (--target-project-id). Moving requires an authorizing move code for the zone, which can be obtained in one of two ways:
  - Generated automatically by this command: omit --code. The configured project is then used to look up the zone and must be its current project. The generated move code is valid for 1 hour.
  - Generated beforehand, by yourself or by whoever controls the zone's current project, via "stackit dns zone generate-move-code" (and optionally checked with "stackit dns zone validate-move-code"), then passed explicitly via --code. The configured project isn't used for the move itself in this case.

```
stackit dns zone move [flags]
```

### Examples

```
  Move a DNS zone with domain name "www.my-zone.com" to project "xxx", automatically generating a move code against the zone's current project
  $ stackit dns zone move --zone-domain-name www.my-zone.com --target-project-id xxx

  Move a DNS zone with domain name "www.my-zone.com" from project "yyy" to project "xxx", automatically generating a move code. --project-id explicitly sets the source project, i.e. the zone's current project
  $ stackit dns zone move --project-id yyy --zone-domain-name www.my-zone.com --target-project-id xxx

  Move a DNS zone with domain name "www.my-zone.com" to project "xxx" using a move code generated beforehand (e.g. via "stackit dns zone generate-move-code")
  $ stackit dns zone move --zone-domain-name www.my-zone.com --target-project-id xxx --code @./move-code.txt
```

### Options

```
      --code secret                Code. Use @<file> to read from a file (example: @./secret.txt), or - to read from stdin (hidden input on a terminal). Literal values are not supported. If not set, a move code is generated automatically for the zone in the configured project, which must then be the zone's current project.
  -h, --help                       Help for "stackit dns zone move"
      --target-project-id string   ID of the project the zone should be moved to
      --zone-domain-name string    Domain name of the zone to move
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

