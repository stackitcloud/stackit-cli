## stackit git instance describe

Shows details of a STACKIT Git instance

### Synopsis

Shows details of a STACKIT Git instance.

```
stackit git instance describe INSTANCE_ID [flags]
```

### Examples

```
  Get details of a STACKIT Git instance with ID "xxx"
  $ stackit git instance describe xxx

  Get details of a STACKIT Git instance with ID "xxx" in JSON format
  $ stackit git instance describe xxx --output-format json
```

### Options

```
  -h, --help[=true|false]   Help for "stackit git instance describe"
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

* [stackit git instance](./stackit_git_instance.md)	 - Provides functionality for STACKIT Git instances

