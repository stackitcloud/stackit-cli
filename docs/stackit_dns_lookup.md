## stackit dns lookup

Looks up a domain on the STACKIT DNS resolvers

### Synopsis

Looks up a domain on the STACKIT DNS resolvers.

```
stackit dns lookup DOMAIN [flags]
```

### Examples

```
  Look up domain "www.example.com"
  $ stackit dns lookup www.example.com

  Look up the AAAA records of domain "www.example.com" on all STACKIT DNS resolvers
  $ stackit dns lookup www.example.com --record-type AAAA --all-resolvers

  Look up domain "www.example.com" including the raw DNS response in JSON format
  $ stackit dns lookup www.example.com --include-raw-response --output-format json
```

### Options

```
      --all-resolvers[=true|false]          If set, all STACKIT DNS resolvers are queried instead of a single one
  -h, --help[=true|false]                   Help for "stackit dns lookup"
      --include-raw-response[=true|false]   If set, the raw DNS response is included in the output. Only shown in JSON and YAML output
      --record-type string                  Record type to look up. If not set, defaults to A (one of: [A, AAAA, SOA, CNAME, NS, MX, TXT, SRV, PTR, ALIAS, DNAME, CAA, DNSKEY, DS, LOC, NAPTR, SSHFP, TLSA, URI, CERT, SVCB, TYPE, CSYNC, HINFO, HTTPS, ANY])
```

### Options inherited from parent commands

```
  -y, --assume-yes[=true|false]   If set, skips all confirmation prompts
      --async[=true|false]        If set, runs the command asynchronously
  -o, --output-format string      Output format, (one of: [json, pretty, none, yaml])
  -p, --project-id string         Project ID
      --region string             Target region for region-specific requests
      --verbosity string          Verbosity of the CLI, (one of: [debug, info, warning, error]) (default "info")
```

### SEE ALSO

* [stackit dns](./stackit_dns.md)	 - Provides functionality for DNS

