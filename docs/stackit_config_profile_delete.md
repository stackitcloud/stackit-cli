## stackit config profile delete

Delete a CLI configuration profile

### Synopsis

Delete a CLI configuration profile.
If the deleted profile is the active profile, the default profile will be set to active.

```
stackit config profile delete PROFILE [flags]
```

### Examples

```
  Delete the configuration profile "my-profile"
  $ stackit config profile delete my-profile
```

### Options

```
  -h, --help   Help for "stackit config profile delete"
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

