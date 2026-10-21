## stackit registry artifactory create

Creates a Container Registry artifactory

### Synopsis

Creates a Container Registry artifactory.

```
stackit registry artifactory create [flags]
```

### Examples

```
  Create an artifactory with name "my-registry"
  $ stackit registry artifactory create --name my-registry
```

### Options

```
  -h, --help          Help for "stackit registry artifactory create"
  -n, --name string   Artifactory name
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

* [stackit registry artifactory](./stackit_registry_artifactory.md)	 - Provides functionality for Container Registry artifactories

