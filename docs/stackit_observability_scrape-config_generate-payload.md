## stackit observability scrape-config generate-payload

Generates a payload to create/update scrape configurations for an Observability instance 

### Synopsis

Generates a JSON payload with values to be used as --payload input for scrape configurations creation or update.
This command can be used to generate a payload to update an existing scrape config or to create a new scrape config job.
To update an existing scrape config job, provide the job name and the instance ID of the Observability instance.
To obtain a default payload to create a new scrape config job, run the command with no flags.
Note that some of the default values provided, such as the job name, the metrics path and URL of the targets, should be adapted to your use case.
See https://docs.api.stackit.cloud/documentation/argus/version/v1#tag/scrape-config/operation/v1_projects_instances_scrapeconfigs_create for information regarding the payload structure.


```
stackit observability scrape-config generate-payload [flags]
```

### Examples

```
  Generate a Create payload with default values, and adapt it with custom values for the different configuration options
  $ stackit observability scrape-config generate-payload --file-path ./payload.json
  <Modify payload in file, if needed>
  $ stackit observability scrape-config create my-config --payload @./payload.json

  Generate an Update payload with the values of an existing configuration named "my-config" for Observability instance xxx, and adapt it with custom values for the different configuration options
  $ stackit observability scrape-config generate-payload --job-name my-config --instance-id xxx --file-path ./payload.json
  <Modify payload in file>
  $ stackit observability scrape-config update my-config --payload @./payload.json

  Generate an Update payload with the values of an existing configuration named "my-config" for Observability instance xxx, and preview it in the terminal
  $ stackit observability scrape-config generate-payload --job-name my-config --instance-id xxx
```

### Options

```
  -f, --file-path string     If set, writes the payload to the given file. If unset, writes the payload to the standard output
  -h, --help                 Help for "stackit observability scrape-config generate-payload"
      --instance-id string   Instance ID
  -n, --job-name string      If set, generates an update payload with the current state of the given scrape config. If unset, generates a create payload with default values
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

* [stackit observability scrape-config](./stackit_observability_scrape-config.md)	 - Provides functionality for scrape configurations in Observability

