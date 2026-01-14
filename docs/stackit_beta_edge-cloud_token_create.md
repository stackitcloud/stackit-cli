## stackit beta edge-cloud token create

Creates a token for an edge instance

### Synopsis

Creates a token for a STACKIT Edge Cloud (STEC) instance.

An expiration time can be set for the token. The expiration time is set in seconds(s), minutes(m), hours(h), days(d) or months(M). Default is 3600 seconds.
Note: the format for the duration is <value><unit>, e.g. 30d for 30 days. You may not combine units.

```
stackit beta edge-cloud token create [flags]
```

### Examples

```
  Create a token for the edge instance with id "xxx".
  $ stackit beta edge-cloud token create --id "xxx"

  Create a token for the edge instance with name "xxx". The token will be valid for one day.
  $ stackit beta edge-cloud token create --name "xxx" --expiration 1d
```

### Options

```
  -e, --expiration string   Expiration time for the kubeconfig, e.g. 5d. By default, the token is valid for 1h.
  -h, --help                Help for "stackit beta edge-cloud token create"
  -i, --id string           The project-unique identifier of this instance.
  -n, --name string         The displayed name to distinguish multiple instances.
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

* [stackit beta edge-cloud token](./stackit_beta_edge-cloud_token.md)	 - Provides functionality for edge service token.

