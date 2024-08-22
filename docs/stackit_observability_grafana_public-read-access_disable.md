## stackit observability grafana public-read-access disable

Disables public read access for Grafana on Observability instances

### Synopsis

Disables public read access for Grafana on Observability instances.
When disabled, a login is required to access the Grafana dashboards of the instance. Otherwise, anyone can access the dashboards.

```
stackit observability grafana public-read-access disable INSTANCE_ID [flags]
```

### Examples

```
  Disable public read access for Grafana on an Observability instance with ID "xxx"
  $ stackit observability grafana public-read-access disable xxx
```

### Options

```
  -h, --help   Help for "stackit observability grafana public-read-access disable"
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit observability grafana public-read-access](./stackit_observability_grafana_public-read-access.md)	 - Enable or disable public read access for Grafana in Observability instances

