## stackit dns record-set describe

Shows details  of a DNS record set

### Synopsis

Shows details  of a DNS record set.

```
stackit dns record-set describe RECORD_SET_ID [flags]
```

### Examples

```
  Get details of DNS record set with ID "xxx" in zone with ID "yyy"
  $ stackit dns record-set describe xxx --zone-id yyy

  Get details of DNS record set with ID "xxx" in zone with ID "yyy" in JSON format
  $ stackit dns record-set describe xxx --zone-id yyy --output-format json
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
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --region string          Target region for region-specific requests
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit dns record-set](./stackit_dns_record-set.md)	 - Provides functionality for DNS record set

