## stackit dns

Provides functionality for DNS

### Synopsis

Provides functionality for DNS

### Examples

```
$ stackit dns zone list --project-id xxx
$ stackit dns zone create --project-id xxx --name my-zone --dns-name my-zone.com
$ stackit dns record-set list --project-id xxx --zone-id xxx
$ stackit dns record-set create --project-id xxx --zone-id xxx --name my-zone --type A --record 1.2.3.4 --record 5.6.7.8
```

### Options

```
  -h, --help   help for dns
```

### Options inherited from parent commands

```
      --project-id string   Project ID
```

### SEE ALSO

* [stackit](./stackit.md)	 - The root command of the STACKIT CLI
* [stackit dns record-set](./stackit_dns_record-set.md)	 - Provides functionality for DNS record set
* [stackit dns zone](./stackit_dns_zone.md)	 - Provides functionality for DNS zone

