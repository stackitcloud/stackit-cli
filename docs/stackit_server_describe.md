## stackit server describe

Shows details of a server

### Synopsis

Shows details of a server.

```
stackit server describe SERVER_ID [flags]
```

### Examples

```
  Show details of a server with ID "xxx"
  $ stackit server describe xxx

  Show details of a server with ID "xxx" in JSON format
  $ stackit server describe xxx --output-format json
```

### Options

```
  -h, --help[=true|false]   Help for "stackit server describe"
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

* [stackit server](./stackit_server.md)	 - Provides functionality for servers

