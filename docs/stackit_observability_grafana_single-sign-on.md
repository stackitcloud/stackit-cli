## stackit observability grafana single-sign-on

Enable or disable single sign-on for Grafana in Observability instances

### Synopsis

Enable or disable single sign-on for Grafana in Observability instances.
When enabled for an instance, overwrites the generic OAuth2 authentication and configures STACKIT single sign-on for that instance.

```
stackit observability grafana single-sign-on [flags]
```

### Options

```
  -h, --help   Help for "stackit observability grafana single-sign-on"
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
* [stackit observability grafana single-sign-on disable](./stackit_observability_grafana_single-sign-on_disable.md)	 - Disables single sign-on for Grafana on Observability instances
* [stackit observability grafana single-sign-on enable](./stackit_observability_grafana_single-sign-on_enable.md)	 - Enables single sign-on for Grafana on Observability instances

