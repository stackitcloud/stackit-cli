## stackit dns flush

Flushes the cache for a domain on the STACKIT DNS resolvers

### Synopsis

Flushes the cache for a domain on the STACKIT DNS resolvers.

```
stackit dns flush DOMAIN [flags]
```

### Examples

```
  Flush the cache for domain "www.example.com"
  $ stackit dns flush www.example.com

  Flush the cache for the A records of domain "www.example.com"
  $ stackit dns flush www.example.com --type A
```

### Options

```
  -h, --help          Help for "stackit dns flush"
      --type string   Record type to flush. If not set, the cache is flushed for all record types of the domain (one of: [A, AAAA, SOA, CNAME, NS, MX, TXT, SRV, PTR, DNAME, CAA, DNSKEY, DS, LOC, NAPTR, SSHFP, TLSA, URI, CERT, SVCB, CSYNC, HINFO, HTTPS])
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

* [stackit dns](./stackit_dns.md)	 - Provides functionality for DNS

