## stackit alpha vpc region create

Creates a regional configuration for a VPC

### Synopsis

Creates a regional configuration for a VPC. VPC can only be used in projects with enabled VPC.

```
stackit alpha vpc region create [flags]
```

### Examples

```
  Create a regional configuration "eu02" for a VPC with ID "xxx"
  $ stackit alpha vpc region create --vpc-id xxx --region eu02

  Create a regional configuration with default DNS name servers, using the configured project and region
  $ stackit alpha vpc region create --vpc-id xxx --ipv4-default-nameservers 8.8.8.8,8.8.4.4
```

### Options

```
  -h, --help                               Help for "stackit alpha vpc region create"
      --ipv4-default-nameservers strings   List of default DNS name server IPs
      --vpc-id string                      VPC ID
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

