## stackit observability grafana single-sign-on disable

Disables single sign-on for Grafana on Observability instances

### Synopsis

Disables single sign-on for Grafana on Observability instances.
When disabled for an instance, the generic OAuth2 authentication is used for that instance.

```
stackit observability grafana single-sign-on disable INSTANCE_ID [flags]
```

### Examples

```
  Disable single sign-on for Grafana on an Observability instance with ID "xxx"
  $ stackit observability grafana single-sign-on disable xxx
```

### Options

```
  -h, --help   Help for "stackit observability grafana single-sign-on disable"
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

* [stackit observability grafana single-sign-on](./stackit_observability_grafana_single-sign-on.md)	 - Enable or disable single sign-on for Grafana in Observability instances

