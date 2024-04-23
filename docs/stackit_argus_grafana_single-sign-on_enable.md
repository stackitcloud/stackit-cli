## stackit argus grafana single-sign-on enable

Enables single sign-on for Grafana on Argus instances

### Synopsis

Enables single sign-on for Grafana on Argus instances.
When enabled for an instance, overwrites the generic OAuth2 authentication and configures STACKIT single sign-on for that instance.

```
stackit argus grafana single-sign-on enable INSTANCE_ID [flags]
```

### Examples

```
  Enable single sign-on for Grafana on an Argus instance with ID "xxx"
  $ stackit argus grafana single-sign-on enable xxx
```

### Options

```
  -h, --help   Help for "stackit argus grafana single-sign-on enable"
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit argus grafana single-sign-on](./stackit_argus_grafana_single-sign-on.md)	 - Enable or disable single sign-on for Grafana in Argus instances

