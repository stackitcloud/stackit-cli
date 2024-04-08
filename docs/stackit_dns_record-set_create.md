## stackit dns record-set create

Creates a DNS record set

### Synopsis

Creates a DNS record set.

```
stackit dns record-set create [flags]
```

### Examples

```
  Create a DNS record set with name "my-rr" with records "1.2.3.4" and "5.6.7.8" in zone with ID "xxx"
  $ stackit dns record-set create --zone-id xxx --name my-rr --record 1.2.3.4 --record 5.6.7.8
```

### Options

```
      --comment string   User comment
  -h, --help             Help for "stackit dns record-set create"
      --name string      Name of the record, should be compliant with RFC1035, Section 2.3.4
      --record strings   Records belonging to the record set
      --ttl int          Time to live, if not provided defaults to the zone's default TTL
      --type string      Record type, one of ["A" "AAAA" "SOA" "CNAME" "NS" "MX" "TXT" "SRV" "PTR" "ALIAS" "DNAME" "CAA"] (default "A")
      --zone-id string   Zone ID
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit dns record-set](./stackit_dns_record-set.md)	 - Provides functionality for DNS record set

