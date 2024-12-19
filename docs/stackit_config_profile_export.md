## stackit config profile export

Exports a CLI configuration profile

### Synopsis

Exports a CLI configuration profile.

```
stackit config profile export PROFILE_NAME [flags]
```

### Examples

```
  Export a profile with name "PROFILE_NAME" to a file in your current directory
  $ stackit config profile export PROFILE_NAME

  Export a profile with name "PROFILE_NAME"" to a specific file path FILE_PATH
  $ stackit config profile export PROFILE_NAME --file-path FILE_PATH
```

### Options

```
  -f, --file-path string   If set, writes the config to the given file path. If unset, writes the config to you current directory with the name of the profile. E.g. '--file-path ~/my-config.json'
  -h, --help               Help for "stackit config profile export"
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

