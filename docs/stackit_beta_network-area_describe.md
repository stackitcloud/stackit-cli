## stackit beta network-area describe

Shows details of a network area

### Synopsis

Shows details of a network area in an organization.

```
stackit beta network-area describe [flags]
```

### Examples

```
  Show details of a network area with ID "xxx" in organization with ID "yyy"
  $ stackit beta network-area describe xxx --organization-id yyy

  Show details of a network area with ID "xxx" in organization with ID "yyy" and show attached projects
  $ stackit beta network-area describe xxx --organization-id yyy --show-attached-projects

  Show details of a network area with ID "xxx" in organization with ID "yyy" in JSON format
  $ stackit beta network-area describe xxx --organization-id yyy --output-format json
```

### Options

```
  -h, --help                     Help for "stackit beta network-area describe"
      --organization-id string   Organization ID
      --show-attached-projects   Whether to show attached projects. If a network area has several attached projects, their retrieval may take some time and the output may be extensive.
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

* [stackit beta network-area](./stackit_beta_network-area.md)	 - Provides functionality for Network Area

