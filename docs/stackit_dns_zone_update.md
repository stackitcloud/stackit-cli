## stackit dns zone update

Updates a DNS zone

### Synopsis

Updates a DNS zone. Performs a partial update; fields not provided are kept unchanged

```
stackit dns zone update [flags]
```

### Examples

```
$ stackit dns zone update --project-id xxx --zone-id xxx --name my-zone --dns-name my-zone.com
```

### Options

```
      --acl string             Access control list
      --contact-email string   Contact email for the zone
      --default-ttl int        Default time to live (default 1000)
      --description string     Description of the zone
      --expire-time int        Expire time
  -h, --help                   help for update
      --name string            User given name of the zone
      --negative-cache int     Negative cache
      --primary strings        Primary name server for secondary zone
      --refresh-time int       Refresh time
      --retry-time int         Retry time
      --zone-id string         Zone ID
```

### Options inherited from parent commands

```
      --project-id string   Project ID
```

### SEE ALSO

* [stackit dns zone](./stackit_dns_zone.md)	 - Provides functionality for DNS zone

