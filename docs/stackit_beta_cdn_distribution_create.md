## stackit beta cdn distribution create

Create a CDN distribution

### Synopsis

Create a CDN distribution for a given originUrl in multiple regions.

```
stackit beta cdn distribution create [flags]
```

### Examples

```
  Create a distribution for regions EU and AF
  $ stackit beta cdn distribution create --regions=EU,AF --origin-url=https://example.com
```

### Options

```
  -h, --help               Help for "stackit beta cdn distribution create"
      --origin-url https   The origin of the content that should be made available through the CDN. Note that the path and query parameters are ignored. Ports are allowed. If no protocol is provided, https is assumed. So `www.example.com:1234/somePath?q=123` is normalized to `https://www.example.com:1234`
      --regions strings    Regions in which content should be cached, multiple of: ["EU" "US" "AF" "SA" "ASIA"] (default [])
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

