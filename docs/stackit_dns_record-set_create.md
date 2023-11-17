## stackit dns record-set create

Creates a DNS record set

### Synopsis

Creates a DNS record set

```
stackit dns record-set create [flags]
```

### Examples

```
$ stackit dns record-set create --project-id xxx --zone-id xxx --name my-zone --type A --record 1.2.3.4 --record 5.6.7.8
```

### Options

```
      --comment string   User comment
  -h, --help             help for create
      --name string      Name of the record, should be compliant with RFC1035, Section 2.3.4
      --record strings   Records belonging to the record set
      --ttl int          Time to live, if not provided defaults to the zone's default TTL
      --type string      Zone type, one of ["A" "AAAA" "SOA" "CNAME" "NS" "MX" "TXT" "SRV" "PTR" "ALIAS" "DNAME" "CAA"]
      --zone-id string   Zone ID
```

### Options inherited from parent commands

```
      --project-id string   Project ID
```

### SEE ALSO

* [stackit dns record-set](./stackit_dns_record-set.md)	 - Provides functionality for DNS record set

