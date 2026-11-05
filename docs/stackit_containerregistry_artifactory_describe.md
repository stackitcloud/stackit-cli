## stackit containerregistry artifactory describe

Shows details of a Container Registry artifactory

### Synopsis

Shows details of a Container Registry artifactory.

```
stackit containerregistry artifactory describe ARTIFACTORY_ID [flags]
```

### Examples

```
  Get details of an artifactory with ID "xxx"
  $ stackit containerregistry artifactory describe xxx

  Get details of an artifactory with ID "xxx" in JSON format
  $ stackit containerregistry artifactory describe xxx --output-format json
```

### Options

```
  -h, --help   Help for "stackit containerregistry artifactory describe"
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

* [stackit containerregistry artifactory](./stackit_containerregistry_artifactory.md)	 - Provides functionality for Container Registry artifactories

