## stackit service-account create

Create a service account

### Synopsis

Create a service account.

```
stackit service-account create [flags]
```

### Examples

```
  Create a service account with name "my-service-account"
  $ stackit service-account create --name my-service-account
```

### Options

```
  -h, --help          Help for "stackit service-account create"
  -n, --name string   Service account name. A unique email will be generated from this name
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit service-account](./stackit_service-account.md)	 - Provides functionality for service accounts

