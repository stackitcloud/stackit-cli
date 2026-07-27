## stackit beta edge-cloud token create

Creates a token for an Edge Cloud instance

### Synopsis

Creates a token for a STACKIT Edge Cloud (STEC) instance.

An expiration time can be set for the token. The expiration time is set in seconds(s), minutes(m), hours(h), days(d) or months(M). Default is 3600(1h) seconds.
Note: the format for the duration is <value><unit>, e.g. 30d for 30 days. You may not combine units.

```
stackit beta edge-cloud token create token for INSTANCE_ID [flags]
```

### Examples

```
  Create a token for the Edge Cloud instance with instance ID "xxx".
  $ stackit beta edge-cloud token create xxx

  Create a token for the Edge Cloud instance with instance ID "xxx". The token will be valid for one day.
  $ stackit beta edge-cloud token create xxx --expiration 1d
```

### Options

```
  -e, --expiration string   Expiration time for the token, e.g. 5d. By default, the token is valid for 1h.
  -h, --help                Help for "stackit beta edge-cloud token create"
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, (one of: [json, pretty, none, yaml])
  -p, --project-id string      Project ID
      --region string          Target region for region-specific requests
      --verbosity string       Verbosity of the CLI, (one of: [debug, info, warning, error]) (default "info")
```

### SEE ALSO

* [stackit beta edge-cloud token](./stackit_beta_edge-cloud_token.md)	 - Provides functionality for Edge Cloud service token.

