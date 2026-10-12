## stackit alpha vpc update

Updates a VPC

### Synopsis

Updates a VPC.

```
stackit alpha vpc update VPC_ID [flags]
```

### Examples

```
  Update a vpc with ID "xxx" to new name "vpc-new"
  $ stackit alpha vpc update xxx --name vpc-new

  Update a vpc with ID "xxx" to new description "updated vpc"
  $ stackit alpha vpc update xxx --description "updated vpc"
```

### Options

```
      --description string      Description of the VPC
  -h, --help                    Help for "stackit alpha vpc update"
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

