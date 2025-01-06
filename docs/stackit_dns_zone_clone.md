## stackit dns zone clone

Clones a DNS zone

### Synopsis

Clones an existing DNS zone with all record sets to a new zone with a different name.

```
stackit dns zone clone ZONE_ID [flags]
```

### Examples

```
  Clones a DNS zone with ID "xxx" to a new zone with DNS name "www.my-zone.com"
  $ stackit dns zone clone xxx --dns-name www.my-zone.com

  Clones a DNS zone with ID "xxx" to a new zone with DNS name "www.my-zone.com" and display name "new-zone"
  $ stackit dns zone clone xxx --dns-name www.my-zone.com --name new-zone

  Clones a DNS zone with ID "xxx" to a new zone with DNS name "www.my-zone.com" and adjust records "true"
  $ stackit dns zone clone xxx --dns-name www.my-zone.com --adjust-records
```

### Options

```
      --adjust-records       Sets content and replaces the DNS name of the original zone with the new DNS name of the cloned zone
      --description string   New description for the cloned zone
      --dns-name string      Fully qualified domain name of the new DNS zone to clone
  -h, --help                 Help for "stackit dns zone clone"
      --name string          User given new name for the cloned zone
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit dns zone](./stackit_dns_zone.md)	 - Provides functionality for DNS zones

