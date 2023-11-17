## stackit dns record-set

Provides functionality for DNS record set

### Synopsis

Provides functionality for DNS record set

### Examples

```
$ stackit dns record-set list --project-id xxx --zone-id xxx
$ stackit dns record-set create --project-id xxx --zone-id xxx --name my-zone --type A --record 1.2.3.4 --record 5.6.7.8
```

### Options

```
  -h, --help   help for record-set
```

### Options inherited from parent commands

```
      --project-id string   Project ID
```

### SEE ALSO

* [stackit dns](./stackit_dns.md)	 - Provides functionality for DNS
* [stackit dns record-set create](./stackit_dns_record-set_create.md)	 - Creates a DNS record set
* [stackit dns record-set delete](./stackit_dns_record-set_delete.md)	 - Delete a DNS record set
* [stackit dns record-set describe](./stackit_dns_record-set_describe.md)	 - Get details of a DNS record set
* [stackit dns record-set list](./stackit_dns_record-set_list.md)	 - List all DNS record sets
* [stackit dns record-set update](./stackit_dns_record-set_update.md)	 - Updates a DNS record set

