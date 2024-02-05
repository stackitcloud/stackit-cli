## stackit dns zone update

Update a DNS zone

### Synopsis

Update a DNS zone.

```
stackit dns zone update ZONE_ID [flags]
```

### Examples

```
  Update the contact email of the DNS zone with ID "xxx"
  $ stackit dns zone update xxx --contact-email someone@domain.com
```

### Options

```
      --acl string             Access control list
      --contact-email string   Contact email for the zone
      --default-ttl int        Default time to live (default 1000)
      --description string     Description of the zone
      --expire-time int        Expire time
  -h, --help                   Help for "stackit dns zone update"
      --name string            User given name of the zone
      --negative-cache int     Negative cache
      --primary strings        Primary name server for secondary zone
      --refresh-time int       Refresh time
      --retry-time int         Retry time
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit dns zone](./stackit_dns_zone.md)	 - Provides functionality for DNS zones

