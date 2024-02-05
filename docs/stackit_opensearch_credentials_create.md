## stackit opensearch credentials create

Create credentials for an OpenSearch instance

### Synopsis

Create credentials (username and password) for an OpenSearch instance.

```
stackit opensearch credentials create [flags]
```

### Examples

```
  Create credentials for an OpenSearch instance
  $ stackit opensearch credentials create --instance-id xxx

  Create credentials for an OpenSearch instance and hide the password in the output
  $ stackit opensearch credentials create --instance-id xxx --hide-password
```

### Options

```
  -h, --help                 Help for "stackit opensearch credentials create"
      --hide-password        Hide password in output
      --instance-id string   Instance ID
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit opensearch credentials](./stackit_opensearch_credentials.md)	 - Provides functionality for OpenSearch credentials

