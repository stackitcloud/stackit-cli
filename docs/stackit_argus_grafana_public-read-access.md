## stackit argus grafana public-read-access

Enable or disable public read access for Grafana in Argus instances

### Synopsis

Enable or disable public read access for Grafana in Argus instances.
When enabled, anyone can access the Grafana dashboards without of the instance logging in. Otherwise, a login is required.

```
stackit argus grafana public-read-access [flags]
```

### Options

```
  -h, --help   Help for "stackit argus grafana public-read-access"
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
* [stackit argus grafana public-read-access disable](./stackit_argus_grafana_public-read-access_disable.md)	 - Disables public read access for Grafana on Argus instances
* [stackit argus grafana public-read-access enable](./stackit_argus_grafana_public-read-access_enable.md)	 - Enables public read access for Grafana on Argus instances

