## stackit config profile import

Imports a CLI configuration profile

### Synopsis

Imports a CLI configuration profile.

```
stackit config profile import [flags]
```

### Examples

```
  Import a config with name "PROFILE_NAME" from file "./config.json"
  $ stackit config profile import --name PROFILE_NAME --config `@./config.json`

  Import a config with name "PROFILE_NAME" from file "./config.json" and do not set as active
  $ stackit config profile import --name PROFILE_NAME --config `@./config.json` --no-set
```

### Options

```
  -c, --config string   File where configuration will be imported from
  -h, --help            Help for "stackit config profile import"
      --name string     Profile name
      --no-set          Set the imported profile not as active
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

