## stackit dns record-set update

Update a DNS record set

### Synopsis

Update a DNS record set.

```
stackit dns record-set update RECORD_SET_ID [flags]
```

### Examples

```
  Update the time to live of the record-set with ID "xxx" for zone with ID "yyy"
  $ stackit dns record-set update xxx --zone-id yyy --ttl 100
```

### Options

```
      --comment string   User comment
  -h, --help             Help for "stackit dns record-set update"
      --name string      Name of the record, should be compliant with RFC1035, Section 2.3.4
      --record strings   Records belonging to the record set. If this flag is used, records already created that aren't set when running the command will be deleted
      --ttl int          Time to live, if not provided defaults to the zone's default TTL
      --zone-id string   Zone ID
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit dns record-set](./stackit_dns_record-set.md)	 - Provides functionality for DNS record set

