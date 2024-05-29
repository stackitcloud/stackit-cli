## stackit config profile list

Lists all CLI configuration profiles

### Synopsis

Lists all CLI configuration profiles.

```
stackit config profile list [flags]
```

### Examples

```
  List the configuration profiles
  $ stackit config profile list

  List the configuration profiles in a json format
  $ stackit config profile list --output-format json
```

### Options

```
  -h, --help   Help for "stackit config profile list"
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit config profile](./stackit_config_profile.md)	 - Manage the CLI configuration profiles

