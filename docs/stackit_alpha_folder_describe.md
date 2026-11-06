## stackit alpha folder describe

Shows details of a STACKIT folder

### Synopsis

Shows details of a STACKIT folder.

```
stackit alpha folder describe [flags]
```

### Examples

```
  Get the details of the configured STACKIT folder
  $ stackit folder describe

  Get the details of a STACKIT folder by explicitly providing the container ID
  $ stackit folder describe xxx

  Get the details of the configured STACKIT folder, including details of the parent resources
  $ stackit folder describe --include-parents
```

### Options

```
  -h, --help              Help for "stackit alpha folder describe"
      --include-parents   When true, the details of the parent resources will be included in the output
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

* [stackit alpha folder](./stackit_alpha_folder.md)	 - Manages the lifecycle of folders.

