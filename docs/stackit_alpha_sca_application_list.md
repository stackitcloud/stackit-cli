## stackit alpha sca application list

Lists all SCA applications

### Synopsis

Lists all STACKIT Container Applications (SCA) applications.

```
stackit alpha sca application list [flags]
```

### Examples

```
  List all SCA applications
  $ stackit alpha sca application list

  List all SCA applications from environment with ID "xxx"
  $ stackit alpha sca application list --environment-id xxx

  List all SCA applications in JSON format
  $ stackit alpha sca application list --output-format json

  List up to 10 SCA applications
  $ stackit alpha sca application list --limit 10
```

### Options

```
      --environment-id string   Optional environment ID to filter applications
  -h, --help                    Help for "stackit alpha sca application list"
      --limit int               Maximum number of entries to list
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

* [stackit alpha sca application](./stackit_alpha_sca_application.md)	 - Provides functionality for SCA applications

