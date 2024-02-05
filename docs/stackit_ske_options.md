## stackit ske options

List SKE provider options

### Synopsis

List STACKIT Kubernetes Engine (SKE) provider options (availability zones, Kubernetes versions, machine images and types, volume types).
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
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit ske](./stackit_ske.md)	 - Provides functionality for SKE

