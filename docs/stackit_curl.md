## stackit curl

Executes an authenticated HTTP request to an endpoint

### Synopsis

Executes an HTTP request to an endpoint, using the authentication provided by the CLI.

```
stackit curl URL [flags]
```

### Examples

```
  Get all the DNS zones for project with ID xxx via GET request to https://dns.api.stackit.cloud/v1/projects/xxx/zones
  $ stackit curl https://dns.api.stackit.cloud/v1/projects/xxx/zones

  Get all the DNS zones for project with ID xxx via GET request to https://dns.api.stackit.cloud/v1/projects/xxx/zones, write complete response (headers and body) to file "./output.txt"
  $ stackit curl https://dns.api.stackit.cloud/v1/projects/xxx/zones -include --output ./output.txt

  Create a new DNS zone for project with ID xxx via POST request to https://dns.api.stackit.cloud/v1/projects/xxx/zones with payload from file "./payload.json"
  $ stackit curl https://dns.api.stackit.cloud/v1/projects/xxx/zones -X POST --data @./payload.json

  Get all the DNS zones for project with ID xxx via GET request to https://dns.api.stackit.cloud/v1/projects/xxx/zones, with header "Authorization: Bearer yyy", fail if server returns error (such as 403 Forbidden)
  $ stackit curl https://dns.api.stackit.cloud/v1/projects/xxx/zones -X POST -H "Authorization: Bearer yyy" --fail
```

### Options

```
      --data string      Content to include in the request body. Can be a string or a file path prefixed with "@"
      --fail             If set, exits with error 22 if response code is 4XX or 5XX
  -H, --header strings   Custom headers to include in the request, can be specified multiple times. If the "Authorization" header is set, it will override the authentication provided by the CLI
  -h, --help             Help for "stackit curl"
      --include          If set, response headers are added to the output
      --output string    Writes output to provided file instead of printing to console
  -X, --request string   HTTP method, defaults to GET
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

* [stackit](./stackit.md)	 - Manage STACKIT resources using the command line

