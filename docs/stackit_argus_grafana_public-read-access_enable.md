## stackit argus grafana public-read-access enable

Enables public read access for Grafana on Argus instances

### Synopsis

Enables public read access for Grafana on Argus instances.
When enabled, anyone can access the Grafana dashboards of the instance without logging in. Otherwise, a login is required.

```
stackit argus grafana public-read-access enable [flags]
```

### Examples

```
  Enable public read access for Grafana on an Argus instance with ID "xxx"
  $ stackit argus grafana public-read-access enable --instance-id xxx
```

### Options

```
  -h, --help                 Help for "stackit argus grafana public-read-access enable"
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

* [stackit argus grafana public-read-access](./stackit_argus_grafana_public-read-access.md)	 - Enable or disable public read access for Grafana in Argus instances

