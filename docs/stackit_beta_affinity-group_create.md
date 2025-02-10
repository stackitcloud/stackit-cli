## stackit beta affinity-group create

Creates an affinity groups

### Synopsis

Creates an affinity groups.

```
stackit beta affinity-group create [flags]
```

### Examples

```
  Create an affinity group with name "AFFINITY_GROUP_NAME" and policy "soft-affinity"
  $ stackit beta affinity-group create --name AFFINITY_GROUP_NAME --policy soft-affinity
```

### Options

```
  -h, --help            Help for "stackit beta affinity-group create"
      --name string     The name of the affinity group.
      --policy string   The policy for the affinity group. Valid values for the policy are: "hard-affinity", "hard-anti-affinity", "soft-affinity", "soft-anti-affinity"
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

* [stackit beta affinity-group](./stackit_beta_affinity-group.md)	 - Manage server affinity groups

