## stackit ske credentials describe

Get details of the credentials associated to a SKE cluster

### Synopsis

Get details of the credentials associated to a STACKIT Kubernetes Engine (SKE) cluster

```
stackit ske credentials describe CLUSTER_NAME [flags]
```

### Examples

```
  Get details of the credentials associated to the SKE cluster with name "my-cluster"
  $ stackit ske credentials describe my-cluster

  Get details of the credentials associated to the SKE cluster with name "my-cluster" in a table format
  $ stackit ske credentials describe my-cluster --output-format pretty
```

### Options

```
  -h, --help   Help for "stackit ske credentials describe"
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit ske credentials](./stackit_ske_credentials.md)	 - Provides functionality for SKE credentials

