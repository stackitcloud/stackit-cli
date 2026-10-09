## stackit alpha vpc describe

Shows details of a VPC

### Synopsis

Shows details of a VPC.

```
stackit alpha vpc describe VPC_ID [flags]
```

### Examples

```
  Show details of a vpc with ID "xxx"
  $ stackit alpha vpc describe xxx

  Show details of a vpc with ID "xxx" in JSON format
  $ stackit alpha vpc describe xxx --output-format json
```

### Options

```
  -h, --help   Help for "stackit alpha vpc describe"
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

