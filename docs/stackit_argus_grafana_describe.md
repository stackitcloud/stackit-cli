## stackit argus grafana describe

Shows details of the Grafana configuration of an Argus instance

### Synopsis

Shows details of the Grafana configuration of an Argus instance.
The Grafana dashboard URL and initial credentials (admin user and password) will be shown in the "pretty" output format. These credentials are only valid for first login. Please change the password after first login. After changing, the initial password is no longer valid.
The initial password is hidden by default, if you want to show it use the "--show-password" flag.

```
stackit argus grafana describe INSTANCE_ID [flags]
```

### Examples

```
  Get details of the Grafana configuration of an Argus instance with ID "xxx"
  $ stackit argus credentials describe xxx

  Get details of the Grafana configuration of an Argus instance with ID "xxx" and show the initial admin password
  $ stackit argus credentials describe xxx --show-password

  Get details of the Grafana configuration of an Argus instance with ID "xxx" in JSON format
  $ stackit argus credentials describe xxx --output-format json
```

### Options

```
  -h, --help            Help for "stackit argus grafana describe"
  -s, --show-password   Show password in output
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

* [stackit argus grafana](./stackit_argus_grafana.md)	 - Provides functionality for the Grafana configuration of Argus instances

