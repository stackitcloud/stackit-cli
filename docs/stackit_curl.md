## stackit curl

Execute an authenticated HTTP request to an endpoint

### Synopsis

Execute an HTTP request to an endpoint, using the authentication provided by the CLI.

```
stackit curl URL [flags]
```

### Examples

```
  Make a GET request to http://locahost:8000
  $ stackit curl http://locahost:8000

  Make a GET request to http://locahost:8000, write complete response (headers and body) to file "./output.txt"
  $ stackit curl http://locahost:8000 -include --output ./output.txt

  Make a POST request to http://locahost:8000 with payload from file "./payload.json"
  $ stackit curl http://locahost:8000 -X POST --data @./payload.json

  Make a POST request to http://locahost:8000 with header "Foo: Bar", fail if server returns error (such as 403 Forbidden)
  $ stackit curl http://locahost:8000 -X POST -H "Foo: Bar" --fail
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
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit](./stackit.md)	 - Manage STACKIT resources using the command line

