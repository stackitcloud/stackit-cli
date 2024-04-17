## stackit argus grafana single-sign-on disable

Disables single sign-on for Grafana on Argus instances

### Synopsis

Disables single sign-on for Grafana on Argus instances.
When disabled for an instance, the generic OAuth2 authentication is used for that instance.

```
stackit argus grafana single-sign-on disable [flags]
```

### Examples

```
  Disable single sign-on for Grafana on an Argus instance with ID "xxx"
  $ stackit argus grafana single-sign-on disable --instance-id xxx
```

### Options

```
  -h, --help                 Help for "stackit argus grafana single-sign-on disable"
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

* [stackit argus grafana single-sign-on](./stackit_argus_grafana_single-sign-on.md)	 - Enable or disable single sign-on for Grafana in Argus instances

