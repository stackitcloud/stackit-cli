## stackit dns record-set delete

Delete a DNS record set

### Synopsis

Delete a DNS record set.

```
stackit dns record-set delete RECORD_SET_ID [flags]
```

### Examples

```
  Delete DNS record set with ID "xxx" in zone with ID "yyy"
  $ stackit dns record-set delete xxx --zone-id yyy
```

### Options

```
  -h, --help             Help for "stackit dns record-set delete"
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

