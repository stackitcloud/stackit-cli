## stackit alpha sca application create

Create a SCA application

### Synopsis

Create a STACKIT Container Applications (SCA) application.

```
stackit alpha sca application create [flags]
```

### Examples

```
  Create a SCA application with name "application-name" and image "my-image" for an environment with ID "yyy"
  $ stackit alpha sca application create --name application-name --image my-image --environment-id yyy

  Create a SCA application with name "application-name" and image "my-image" with 2 instances
  $ stackit alpha sca application create --name application-name --image my-image --instances 2

  Create a SCA application with name "application-name" and image "my-image" exposing port 8888 of the container
  $ stackit alpha sca application create --name application-name --image my-image --external-port 8888

  Create a SCA application with name "application-name" and image "my-image" disabling public networking
  $ stackit alpha sca application create --name application-name --image my-image --public=false

  Create a SCA application with name "application-name" and image "my-image" and environment variables ENV1=value1 and ENV2=value2
  $ stackit alpha sca application create --name application-name --image my-image --environment-vars ENV1=value1,ENV2=value2
```

### Options

```
      --args strings                      Arguments to pass to the application container command
      --commands strings                  Commands to execute in the application container
      --container-name string             Container name (default "container-1")
      --cpu int32                         The dedicated virtual CPU processing power allocated per container instance (default 1000)
      --environment-id string             Environment ID (uses default environment if not set)
      --environment-vars stringToString   Environment variables to inject into the application (default [])
      --external-port int32               Container external exposed port (default 8080)
  -h, --help                              Help for "stackit alpha sca application create"
      --http-rule-concurrency int32       Target number of in-flight requests to trigger autoscaling (if autoscaling is enabled)
      --http-rule-rps int32               Target number of requests per second to trigger autoscaling (if autoscaling is enabled)
      --image string                      Container image
      --instances int32                   The number of application instances (if manually scaled) (default 1)
      --max-instances int32               The maximum number of application instances (if autoscaling is enabled)
      --memory int32                      The total amount of memory (RAM) allocated per container instance (default 1024)
      --min-instances int32               The minimum number of application instances (if autoscaling is enabled) (default 1)
      --name string                       Application display name
      --public                            Exposes your application securely to the public internet via HTTPS endpoint (default true)
      --scale-to-zero                     Enable scale to zero (if autoscaling is enabled)
      --scaling-type string               Scaling type, (one of: [manual, auto]) (default "manual")
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

