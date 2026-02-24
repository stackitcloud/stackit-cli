## stackit organization describe

Show an organization

### Synopsis

Show an organization.

```
stackit organization describe [flags]
```

### Examples

```
  Describe the organization with the organization uuid "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  $ stackit organization describe xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx

  Describe the organization with the container id "foo-bar-organization"
  $ stackit organization describe foo-bar-organization
```

### Options

```
  -h, --help   Help for "stackit organization describe"
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

* [stackit organization](./stackit_organization.md)	 - Manages organizations

