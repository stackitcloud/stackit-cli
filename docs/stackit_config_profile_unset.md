## stackit config profile unset

Unset the current active CLI configuration profile

### Synopsis

Unset the current active CLI configuration profile.
When no profile is set, the default profile will be used.

```
stackit config profile unset [flags]
```

### Examples

```
  Unset the currently active configuration profile. The default profile will be used.
  $ stackit config profile unset
```

### Options

```
  -h, --help   Help for "stackit config profile unset"
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

* [stackit config profile](./stackit_config_profile.md)	 - Manage the CLI configuration profiles

