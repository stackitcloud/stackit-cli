## stackit observability grafana public-read-access enable

Enables public read access for Grafana on Observability instances

### Synopsis

Enables public read access for Grafana on Observability instances.
When enabled, anyone can access the Grafana dashboards of the instance without logging in. Otherwise, a login is required.

```
stackit observability grafana public-read-access enable INSTANCE_ID [flags]
```

### Examples

```
  Enable public read access for Grafana on an Observability instance with ID "xxx"
  $ stackit observability grafana public-read-access enable xxx
```

### Options

```
  -h, --help   Help for "stackit observability grafana public-read-access enable"
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

* [stackit observability grafana public-read-access](./stackit_observability_grafana_public-read-access.md)	 - Enable or disable public read access for Grafana in Observability instances

