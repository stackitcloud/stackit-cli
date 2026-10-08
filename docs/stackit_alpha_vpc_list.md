## stackit alpha vpc list

Lists all VPC of a project

### Synopsis

Lists all VPC of a project.

```
stackit alpha vpc list [flags]
```

### Examples

```
  Lists all VPCs
  $ stackit alpha vpc list xxx

  Lists all VPCsin JSON format
  $ stackit alpha vpc list --output-format json

  Lists up to 10 VPCs
  $ stackit alpha vpc list xxx --limit 10

  Lists all VPCs which has the name "my-vpc"
  $ stackit alpha vpc list xxx --filter "name == 'my-vpc'"
```

### Options

```
      --filter string   Filter resources by fields. A subset of expr-lang is supported. See https://expr-lang.org/docs/language-definition for usage details
  -h, --help            Help for "stackit alpha vpc list"
      --limit int       Maximum number of VPCs to return
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

* [stackit alpha vpc](./stackit_alpha_vpc.md)	 - Manages vpcs

