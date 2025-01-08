## stackit beta server volume describe

Describes a server volume attachment

### Synopsis

Describes a server volume attachment.

```
stackit beta server volume describe VOLUME_ID [flags]
```

### Examples

```
  Get details of the attachment of volume with ID "xxx" to server with ID "yyy"
  $ stackit beta server volume describe xxx --server-id yyy

  Get details of the attachment of volume with ID "xxx" to server with ID "yyy" in JSON format
  $ stackit beta server volume describe xxx --server-id yyy --output-format json

  Get details of the attachment of volume with ID "xxx" to server with ID "yyy" in yaml format
  $ stackit beta server volume describe xxx --server-id yyy --output-format yaml
```

### Options

```
  -h, --help               Help for "stackit beta server volume describe"
      --server-id string   Server ID
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

* [stackit beta server volume](./stackit_beta_server_volume.md)	 - Provides functionality for server volumes

