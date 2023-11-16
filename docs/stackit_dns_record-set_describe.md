## stackit dns record-set describe

Get details of a DNS record set

### Synopsis

Get details of a DNS record set

```
stackit dns record-set describe RECORD_SET_ID [flags]
```

### Examples

```
  Get details of DNS record set with ID "xxx" in zone with ID "yyy"
  $ stackit dns record-set describe xxx --zone-id yyy

  Get details of DNS record set with ID "xxx" in zone with ID "yyy" in a table format
  $ stackit dns record-set describe xxx --zone-id yyy --output-format pretty
```

### Options

```
  -h, --help             Help for "stackit dns record-set describe"
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

