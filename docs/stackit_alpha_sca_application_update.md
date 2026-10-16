## stackit alpha sca application update

Update a SCA application

### Synopsis

Update a STACKIT Container Applications (SCA) application.

```
stackit alpha sca application update [flags]
```

### Examples

```
  Update the number of instances of a SCA application with ID "xxx" from an environment with ID "yyy"
  $ stackit alpha sca application update xxx --instances 2 --environment-id yyy

  Update the container image of a SCA application with ID "xxx" from an environment with ID "yyy"
  $ stackit alpha sca application update xxx --image new-image --environment-id yyy
```

### Options

```
      --args strings                      Arguments to pass to the application container command
      --commands strings                  Commands to execute in the application container
      --cpu int32                         The dedicated virtual CPU processing power allocated per container instance
      --environment-id string             Environment ID (uses default environment if not set)
      --environment-vars stringToString   Environment variables to inject into the application (default [])
      --external-port int32               Container external exposed port
  -h, --help                              Help for "stackit alpha sca application update"
      --http-rule-concurrency int32       Target number of in-flight requests to trigger autoscaling (if autoscaling is enabled)
      --http-rule-rps int32               Target number of requests per second to trigger autoscaling (if autoscaling is enabled)
      --image string                      Container image
      --instances int32                   The number of application instances (if manually scaled)
      --max-instances int32               The maximum number of application instances (if autoscaling is enabled)
      --memory int32                      The total amount of memory (RAM) allocated per container instance
      --min-instances int32               The minimum number of application instances (if autoscaling is enabled)
      --public                            Exposes your application securely to the public internet via HTTPS endpoint
      --scale-to-zero                     Enable scale to zero (if autoscaling is enabled)
      --scaling-type string               Scaling type (one of: [manual, auto])
      --stopped                           Stopped
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, (one of: [json, pretty, none, yaml])
  -p, --project-id string      Project ID
      --region string          Target region for region-specific requests
      --verbosity string       Verbosity of the CLI, (one of: [debug, info, warning, error]) (default "info")
```

### SEE ALSO

* [stackit alpha sca application](./stackit_alpha_sca_application.md)	 - Provides functionality for SCA applications

