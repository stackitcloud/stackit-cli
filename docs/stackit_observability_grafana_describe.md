## stackit observability grafana describe

Shows details of the Grafana configuration of an Observability instance

### Synopsis

Shows details of the Grafana configuration of an Observability instance.
The Grafana dashboard URL and initial credentials (admin user and password) will be shown in the "pretty" output format. These credentials are only valid for first login. Please change the password after first login. After changing, the initial password is no longer valid.
The initial password is hidden by default, if you want to show it use the "--show-password" flag.

```
stackit observability grafana describe INSTANCE_ID [flags]
```

### Examples

```
  Get details of the Grafana configuration of an Observability instance with ID "xxx"
  $ stackit observability grafana describe xxx

  Get details of the Grafana configuration of an Observability instance with ID "xxx" and show the initial admin password
  $ stackit observability grafana describe xxx --show-password

  Get details of the Grafana configuration of an Observability instance with ID "xxx" in JSON format
  $ stackit observability grafana describe xxx --output-format json
```

### Options

```
  -h, --help            Help for "stackit observability grafana describe"
  -s, --show-password   Show password in output
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

* [stackit observability grafana](./stackit_observability_grafana.md)	 - Provides functionality for the Grafana configuration of Observability instances

