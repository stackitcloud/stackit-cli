## stackit registry artifactory update

Updates a Container Registry artifactory

### Synopsis

Updates a Container Registry artifactory.

```
stackit registry artifactory update ARTIFACTORY_ID [flags]
```

### Examples

```
  Update the proxy configuration for an artifactory with ID "xxx"
  $ stackit registry artifactory update xxx --proxy-url https://registry-1.docker.io
```

### Options

```
  -h, --help                         Help for "stackit registry artifactory update"
      --proxy-access-key string      Access key for upstream proxy registry authentication
      --proxy-access-secret string   Proxy-Access-Secret. Can be a string (deprecated) or a file path, if prefixed with '@' (example: @./secret.txt). Will be read from stdin when empty.
      --proxy-description string     Description of the proxy registry
      --proxy-url string             Upstream registry URL to proxy
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

