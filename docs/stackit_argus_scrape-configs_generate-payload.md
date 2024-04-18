## stackit argus scrape-configs generate-payload

Generates a payload to create/update Scrape Configurations for an Argus instance 

### Synopsis

Generates a JSON payload with values to be used as --payload input for Scrape Configurations creation or update.
This command can be used to generate a payload to update an existing Scrape Config job or to create a new Scrape Config job.
To update an existing Scrape Config job, provide the job name and the instance ID of the Argus instance.
To obtain a default payload to create a new Scrape Config job, run the command with no flags.
See https://docs.api.stackit.cloud/documentation/argus/version/v1#tag/scrape-config/operation/v1_projects_instances_scrapeconfigs_create for information regarding the payload structure.


```
stackit argus scrape-configs generate-payload [flags]
```

### Examples

```
  Generate a Create payload with default values, and adapt it with custom values for the different configuration options
  $ stackit argus scrape-configs generate-payload > ./payload.json
  <Modify payload in file, if needed>
  $ stackit argus scrape-configs create my-config --payload @./payload.json

  Generate an Update payload with the values of an existing configuration named "my-config" for Argus instance xxx, and adapt it with custom values for the different configuration options
  $ stackit argus scrape-configs generate-payload --job-name my-config --instance-id xxx > ./payload.json
  <Modify payload in file>
  $ stackit argus scrape-configs update my-config --payload @./payload.json
```

### Options

```
  -h, --help                 Help for "stackit argus scrape-configs generate-payload"
      --instance-id string   Instance ID
  -n, --job-name string      If set, generates an update payload with the current state of the given scrape config. If unset, generates a create payload with default values
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

* [stackit argus scrape-configs](./stackit_argus_scrape-configs.md)	 - Provides functionality for scraping configs in Argus.

