## stackit alpha vpc region update

Updates a regional configuration for a VPC

### Synopsis

Updates a regional configuration for a VPC.

```
stackit alpha vpc region update [flags]
```

### Examples

```
  Update the default DNS name servers for a VPC with ID "xxx" in region "eu02"
  $ stackit alpha vpc region update --vpc-id xxx --region eu02 --ipv4-default-nameservers 8.8.8.8,8.8.4.4

  Clear the default DNS name servers, using the configured project and region
  $ stackit alpha vpc region update --vpc-id xxx --ipv4-default-nameservers ""
```

### Options

```
  -h, --help                               Help for "stackit alpha vpc region update"
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

