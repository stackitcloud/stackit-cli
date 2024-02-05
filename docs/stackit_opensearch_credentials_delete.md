## stackit opensearch credentials delete

Delete credentials of an OpenSearch instance

### Synopsis

Delete credentials of an OpenSearch instance.

```
stackit opensearch credentials delete CREDENTIALS_ID [flags]
```

### Examples

```
  Delete credentials with ID "xxx" of OpenSearch instance with ID "yyy"
  $ stackit opensearch credentials delete xxx --instance-id yyy
```

### Options

```
  -h, --help                 Help for "stackit opensearch credentials delete"
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

