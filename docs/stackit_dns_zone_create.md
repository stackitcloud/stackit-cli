## stackit dns zone create

Creates a DNS zone

### Synopsis

Creates a DNS zone

```
stackit dns zone create [flags]
```

### Examples

```
$ stackit dns zone create --project-id xxx --name my-zone --dns-name my-zone.com
```

### Options

```
      --acl string             Access control list
      --contact-email string   Contact email for the zone
      --default-ttl int        Default time to live (default 1000)
      --description string     Description of the zone
      --dns-name string        DNS zone name
      --expire-time int        Expire time
  -h, --help                   help for create
      --is-reverse-zone        Is reverse zone
      --name string            User given name of the zone
      --negative-cache int     Negative cache
      --primary strings        Primary name server for secondary zone
      --refresh-time int       Refresh time
      --retry-time int         Retry time
      --type string            Zone type
```

### Options inherited from parent commands

```
      --project-id string   Project ID
```

### SEE ALSO

* [stackit dns zone](./stackit_dns_zone.md)	 - Provides functionality for DNS zone

