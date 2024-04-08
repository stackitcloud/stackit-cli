## stackit dns zone create

Creates a DNS zone

### Synopsis

Creates a DNS zone.

```
stackit dns zone create [flags]
```

### Examples

```
  Create a DNS zone with name "my-zone" and DNS name "www.my-zone.com"
  $ stackit dns zone create --name my-zone --dns-name www.my-zone.com

  Create a DNS zone with name "my-zone", DNS name "www.my-zone.com" and default time to live of 1000ms
  $ stackit dns zone create --name my-zone --dns-name www.my-zone.com --default-ttl 1000
```

### Options

```
      --acl string             Access control list
      --contact-email string   Contact email for the zone
      --default-ttl int        Default time to live (default 1000)
      --description string     Description of the zone
      --dns-name string        Fully qualified domain name of the DNS zone
      --expire-time int        Expire time
  -h, --help                   Help for "stackit dns zone create"
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
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit dns zone](./stackit_dns_zone.md)	 - Provides functionality for DNS zones

