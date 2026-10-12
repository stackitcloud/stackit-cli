## stackit alpha vpc create

Creates a VPC

### Synopsis

Creates a VPC.
VPC can only be used in projects with enabled VPC.

```
stackit alpha vpc create [flags]
```

### Examples

```
  Create a vpc with name "my-vpc"
  $ stackit alpha vpc create --name my-vpc
```

### Options

```
      --description string      Description of the VPC
  -h, --help                    Help for "stackit alpha vpc create"
      --labels stringToString   Comma separated list of labels of the VPC (default [])
      --name string             Name of the VPC
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

