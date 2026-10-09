## stackit alpha vpc region delete

Deletes a regional configuration for a VPC

### Synopsis

Deletes a regional configuration for a VPC.

```
stackit alpha vpc region delete [flags]
```

### Examples

```
  Delete the regional configuration "eu02" for a VPC with ID "xxx"
  $ stackit alpha vpc region delete --vpc-id xxx --region eu02

  Delete a regional configuration using the configured project and region
  $ stackit alpha vpc region delete --vpc-id xxx
```

### Options

```
  -h, --help            Help for "stackit alpha vpc region delete"
      --vpc-id string   VPC ID
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

* [stackit alpha vpc region](./stackit_alpha_vpc_region.md)	 - Manages regional configurations of a VPC

