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

  Create a proxy-cache artifactory for Docker Hub
  $ stackit registry artifactory create --name docker-proxy --proxy-url https://registry-1.docker.io --proxy-registry-type docker-hub

  Create a proxy-cache artifactory with credentials
  $ stackit registry artifactory create --name my-proxy --proxy-url https://custom-registry.com --proxy-registry-type harbor --proxy-access-key my-user --proxy-access-secret my-secret
```

### Options

```
  -h, --help                         Help for "stackit registry artifactory create"
  -n, --name string                  Artifactory name
      --proxy-access-key string      Access key for upstream proxy registry authentication
      --proxy-access-secret string   Access secret for upstream proxy registry authentication
      --proxy-description string     Description of the proxy registry
      --proxy-registry-type string   Registry type of the upstream proxy (e.g., docker-hub, harbor)
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

