## stackit dns zone delete

Deletes a DNS zone

### Synopsis

Deletes a DNS zone.

```
stackit dns zone delete ZONE_ID [flags]
```

### Examples

```
  Delete a DNS zone with ID "xxx"
  $ stackit dns zone delete xxx
```

### Options

```
  -h, --help   Help for "stackit dns zone delete"
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit dns zone](./stackit_dns_zone.md)	 - Provides functionality for DNS zones

