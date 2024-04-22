## stackit argus scrape-configs create

Creates a Scrape Config Job for an Argus instance

### Synopsis

Creates a Scrape Config Job for an Argus instance.
The payload can be provided as a JSON string or a file path prefixed with "@".
If no payload is provided, a default payload will be used.
See https://docs.api.stackit.cloud/documentation/argus/version/v1#tag/scrape-config/operation/v1_projects_instances_scrapeconfigs_create for information regarding the payload structure.

```
stackit argus scrape-configs create [flags]
```

### Examples

```
  Create a Scrape Config job using default configuration
  $ stackit argus scrape-configs create

  Create a Scrape Config job using an API payload sourced from the file "./payload.json"
  $ stackit argus scrape-configs create --payload @./payload.json

  Create a Scrape Config job using an API payload provided as a JSON string
  $ stackit argus scrape-configs create --payload "{...}"

  Generate a payload with default values, and adapt it with custom values for the different configuration options
  $ stackit argus scrape-configs generate-payload > ./payload.json
  <Modify payload in file, if needed>
  $ stackit argus scrape-configs create --payload @./payload.json
```

### Options

```
  -h, --help                 Help for "stackit argus scrape-configs create"
      --instance-id string   Instance ID
      --payload string       Request payload (JSON). Can be a string or a file path, if prefixed with "@" (example: @./payload.json). If unset, will use a default payload (you can check it by running "stackit argus scrape-configs generate-payload")
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit argus scrape-configs](./stackit_argus_scrape-configs.md)	 - Provides functionality for scrape configs in Argus.

