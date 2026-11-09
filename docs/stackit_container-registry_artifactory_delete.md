## stackit container-registry artifactory delete

Deletes a Container Registry artifactory

### Synopsis

Deletes a Container Registry artifactory.

```
stackit container-registry artifactory delete ARTIFACTORY_ID [flags]
```

### Examples

```
  Delete an artifactory with ID "xxx"
  $ stackit container-registry artifactory delete xxx
```

### Options

```
  -h, --help   Help for "stackit container-registry artifactory delete"
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

* [stackit container-registry artifactory](./stackit_container-registry_artifactory.md)	 - Provides functionality for Container Registry artifactories

