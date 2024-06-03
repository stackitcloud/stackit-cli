## stackit config profile set

Set a CLI configuration profile

### Synopsis

Set a CLI configuration profile as the active profile.
The profile to be used can be managed via the STACKIT_CLI_PROFILE environment variable or using the "stackit config profile set PROFILE" and "stackit config profile unset" commands.
The environment variable takes precedence over what is set via the commands.
When no profile is set, the default profile is used.

```
stackit config profile set PROFILE [flags]
```

### Examples

```
  Set the configuration profile "my-profile" as the active profile
  $ stackit config profile set my-profile
```

### Options

```
  -h, --help   Help for "stackit config profile set"
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

