## stackit container-registry artifactory update

Updates a Container Registry artifactory

### Synopsis

Updates a Container Registry artifactory.

```
stackit container-registry artifactory update ARTIFACTORY_ID [flags]
```

### Examples

```
  Update the upstream proxy URL for an artifactory with ID "xxx"
  $ stackit container-registry artifactory update xxx --proxy-url https://registry-1.docker.io

  Update the proxy description for an artifactory with ID "xxx"
  $ stackit container-registry artifactory update xxx --proxy-description "Updated proxy description"

  Update the proxy credentials reading secret from stdin for an artifactory with ID "xxx"
  $ stackit container-registry artifactory update xxx --proxy-access-key my-user --proxy-access-secret -

  Update the proxy credentials providing path to a file with secret for an artifactory with ID "xxx"
  $ stackit container-registry artifactory update xxx --proxy-access-key my-user --proxy-access-secret @./secret.txt

  Update all proxy configuration fields for an artifactory with ID "xxx"
  $ stackit container-registry artifactory update xxx --proxy-url https://registry-1.docker.io --proxy-description "Updated mirror" --proxy-access-key my-user --proxy-access-secret @./secret.txt
```

### Options

```
  -h, --help[=true|false]            Help for "stackit container-registry artifactory update"
      --proxy-access-key string      Access key for upstream proxy registry authentication
      --proxy-access-secret secret   Proxy-Access-Secret. Use @<file> to read from a file (example: @./secret.txt), or - to read from stdin (hidden input on a terminal). Literal values are not supported.
      --proxy-description string     Description of the proxy registry
      --proxy-url string             Upstream registry URL to proxy
```

### Options inherited from parent commands

```
  -y, --assume-yes[=true|false]   If set, skips all confirmation prompts
      --async[=true|false]        If set, runs the command asynchronously
  -o, --output-format string      Output format, (one of: [json, pretty, none, yaml])
  -p, --project-id string         Project ID
      --region string             Target region for region-specific requests
      --verbosity string          Verbosity of the CLI, (one of: [debug, info, warning, error]) (default "info")
```

### SEE ALSO

* [stackit container-registry artifactory](./stackit_container-registry_artifactory.md)	 - Provides functionality for Container Registry artifactories

