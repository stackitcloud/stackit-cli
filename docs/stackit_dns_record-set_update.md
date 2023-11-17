## stackit dns record-set update

Updates a DNS record set

### Synopsis

Updates a DNS record set. Performs a partial update; fields not provided are kept unchanged

```
stackit dns record-set update [flags]
```

### Examples

```
$ stackit dns record-set update --project-id xxx --zone-id xxx --record-set-id xxx --name my-zone --type A --record 1.2.3.4 --record 5.6.7.8
```

### Options

```
      --comment string         User comment
  -h, --help                   help for update
      --name string            Name of the record, should be compliant with RFC1035, Section 2.3.4
      --record strings         Records belonging to the record set. If this flag is used, records already created that aren't set when running the command will be deleted
      --record-set-id string   Record set ID
      --ttl int                Time to live, if not provided defaults to the zone's default TTL
      --zone-id string         Zone ID
```

### Options inherited from parent commands

```
      --project-id string   Project ID
```

### SEE ALSO

* [stackit dns record-set](./stackit_dns_record-set.md)	 - Provides functionality for DNS record set

