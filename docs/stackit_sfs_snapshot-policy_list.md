## stackit sfs snapshot-policy list

Lists all snapshot policies of a project

### Synopsis

Lists all snapshot policies of a project.

```
stackit sfs snapshot-policy list [flags]
```

### Examples

```
  List all snapshot policies
  $ stackit sfs snapshot-policy list

  List only mutable snapshot policies
  $ stackit sfs snapshot-policy list --immutable mutable-only

  List only immutable snapshot policies
  $ stackit sfs snapshot-policy list --immutable immutable-only

  List up to 10 snapshot policies
  $ stackit sfs snapshot-policy list --limit 10
```

### Options

```
  -h, --help               Help for "stackit sfs snapshot-policy list"
      --immutable string   Immutable snapshot policy, (one of: [all, immutable-only, mutable-only]) (default "all")
      --limit int          Maximum number of entries to list
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

* [stackit sfs snapshot-policy](./stackit_sfs_snapshot-policy.md)	 - Provides functionality for SFS snapshot policies

