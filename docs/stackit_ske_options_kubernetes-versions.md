## stackit ske options kubernetes-versions

Lists SKE provider options for kubernetes-versions

### Synopsis

Lists STACKIT Kubernetes Engine (SKE) provider options for kubernetes-versions.

```
stackit ske options kubernetes-versions [flags]
```

### Examples

```
  List SKE options for kubernetes-versions
  $ stackit ske options kubernetes-versions

  List SKE options for supported kubernetes-versions
  $ stackit ske options kubernetes-versions --supported
```

### Options

```
  -h, --help        Help for "stackit ske options kubernetes-versions"
      --supported   List supported versions only
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

* [stackit ske options](./stackit_ske_options.md)	 - Lists SKE provider options

