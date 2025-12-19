## stackit beta cdn distribution list

List CDN distributions

### Synopsis

List all CDN distributions in your account.

```
stackit beta cdn distribution list [flags]
```

### Examples

```
  List all CDN distributions
  $ stackit beta cdn distribution list

  List all CDN distributions sorted by id
  $ stackit beta cdn distribution list --sort-by=id
```

### Options

```
      -- int             Limit the output to the first n elements
  -h, --help             Help for "stackit beta cdn distribution list"
      --sort-by string   Sort entries by a specific field, one of ["id" "createdAt" "updatedAt" "originUrl" "status" "originUrlRelated"] (default "createdAt")
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

* [stackit beta cdn distribution](./stackit_beta_cdn_distribution.md)	 - Manage CDN distributions

