## stackit argus grafana describe

Shows details of the Grafana configuration of an Argus instance

### Synopsis

Shows details of the Grafana configuration of an Argus instance.

```
stackit argus grafana describe [flags]
```

### Examples

```
  Get details of the Grafana configuration of an Argus instance with ID "xxx"
  $ stackit argus credentials describe --instance-id xxx

  Get details of of the Grafana configuration of an Argus instance with ID "xxx" in a table format
  $ stackit argus credentials describe --instance-id xxx --output-format pretty
```

### Options

```
  -h, --help                 Help for "stackit argus grafana describe"
      --instance-id string   Instance ID
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit argus grafana](./stackit_argus_grafana.md)	 - Provides functionality for the Grafana configuration of Argus instances

