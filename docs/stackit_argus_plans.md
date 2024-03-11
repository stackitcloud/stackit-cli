## stackit argus plans

Lists all Argus service plans

### Synopsis

Lists all Argus service plans.

```
stackit argus plans [flags]
```

### Examples

```
  List all Argus service plans
  $ stackit argus plans

  List all Argus service plans in JSON format
  $ stackit argus plans --output-format json

  List up to 10 Argus service plans
  $ stackit argus plans --limit 10
```

### Options

```
  -h, --help        Help for "stackit argus plans"
      --limit int   Maximum number of entries to list
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit argus](./stackit_argus.md)	 - Provides functionality for Argus

