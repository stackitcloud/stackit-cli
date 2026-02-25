## stackit ske options

Lists SKE provider options

### Synopsis

Lists STACKIT Kubernetes Engine (SKE) provider options (availability zones, Kubernetes versions, machine images and types, volume types).
Pass one or more flags to filter what categories are shown.

```
stackit ske options [flags]
```

### Examples

```
  List SKE options for all categories
  $ stackit ske options

  List SKE options regarding Kubernetes versions only
  $ stackit ske options --kubernetes-versions

  List SKE options regarding Kubernetes versions and machine images
  $ stackit ske options --kubernetes-versions --machine-images
```

### Options

```
      --availability-zones    Lists availability zones
  -h, --help                  Help for "stackit ske options"
      --kubernetes-versions   Lists supported kubernetes versions
      --machine-images        Lists supported machine images
      --machine-types         Lists supported machine types
      --volume-types          Lists supported volume types
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

* [stackit ske](./stackit_ske.md)	 - Provides functionality for SKE
* [stackit ske options availability-zones](./stackit_ske_options_availability-zones.md)	 - Lists SKE provider options for availability-zones
* [stackit ske options kubernetes-versions](./stackit_ske_options_kubernetes-versions.md)	 - Lists SKE provider options for kubernetes-versions
* [stackit ske options machine-images](./stackit_ske_options_machine-images.md)	 - Lists SKE provider options for machine-images
* [stackit ske options machine-types](./stackit_ske_options_machine-types.md)	 - Lists SKE provider options for machine-types
* [stackit ske options volume-types](./stackit_ske_options_volume-types.md)	 - Lists SKE provider options for volume-types

