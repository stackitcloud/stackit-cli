## stackit observability grafana public-read-access

Enable or disable public read access for Grafana in Observability instances

### Synopsis

Enable or disable public read access for Grafana in Observability instances.
When enabled, anyone can access the Grafana dashboards of the instance without logging in. Otherwise, a login is required.

```
stackit observability grafana public-read-access [flags]
```

### Options

```
  -h, --help   Help for "stackit observability grafana public-read-access"
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

* [stackit observability grafana](./stackit_observability_grafana.md)	 - Provides functionality for the Grafana configuration of Observability instances
* [stackit observability grafana public-read-access disable](./stackit_observability_grafana_public-read-access_disable.md)	 - Disables public read access for Grafana on Observability instances
* [stackit observability grafana public-read-access enable](./stackit_observability_grafana_public-read-access_enable.md)	 - Enables public read access for Grafana on Observability instances

