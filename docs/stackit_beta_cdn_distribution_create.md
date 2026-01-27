## stackit beta cdn distribution create

Create a CDN distribution

### Synopsis

Create a CDN distribution for a given originUrl in multiple regions.

```
stackit beta cdn distribution create [flags]
```

### Examples

```
  Create a CDN distribution with an HTTP backend
  $ stackit beta cdn distribution create --http  --http-origin-url https://example.com \
--regions AF,EU

  Create a CDN distribution with an Object Storage backend
  $ stackit beta cdn distribution create --bucket --bucket-url https://bucket.example.com \
--bucket-credentials-access-key-id yyyy --bucket-region EU \
--regions AF,EU

  Create a CDN distribution passing the password via stdin, take care that there's a '\n' at the end of the input'
  $ cat secret.txt | stackit beta cdn distribution create -y --project-id xxx \
--bucket --bucket-url https://bucket.example.com --bucekt-credentials-access-key-id yyyy --bucket-region EU \
--regions AF,EU
```

### Options

```
      --blocked-countries strings                 Comma-separated list of ISO 3166-1 alpha-2 country codes to block (e.g., 'US,DE,FR')
      --blocked-ips strings                       Comma-separated list of IPv4 addresses to block (e.g., '10.0.0.8,127.0.0.1')
      --bucket                                    Use Object Storage backend
      --bucket-credentials-access-key-id string   Access Key ID for Object Storage backend
      --bucket-region string                      Region for Object Storage backend
      --bucket-url string                         Bucket URL for Object Storage backend
      --default-cache-duration string             ISO8601 duration string for default cache duration (e.g., 'PT1H30M' for 1 hour and 30 minutes)
  -h, --help                                      Help for "stackit beta cdn distribution create"
      --http                                      Use HTTP backend
      --http-geofencing stringArray               Geofencing rules for HTTP backend in the format 'https://example.com US,DE'. URL and countries have to be quoted. Repeatable.
      --http-origin-request-headers strings       Origin request headers for HTTP backend in the format 'HeaderName: HeaderValue', repeatable. WARNING: do not store sensitive values in the headers!
      --http-origin-url string                    Origin URL for HTTP backend
      --loki                                      Enable Loki log sink for the CDN distribution
      --loki-push-url string                      Push URL for log sink
      --loki-username string                      Username for log sink
      --monthly-limit-bytes int                   Monthly limit in bytes for the CDN distribution
      --optimizer                                 Enable optimizer for the CDN distribution (paid feature).
      --regions strings                           Regions in which content should be cached, multiple of: ["EU" "US" "AF" "SA" "ASIA"] (default [])
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

