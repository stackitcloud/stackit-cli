## stackit beta routing-table

Manage routing-tables and its according routes

### Synopsis

Manage routing tables and their associated routes.

This functionality is currently in BETA. At this stage, only listing and describing
routing-tables, as well as full CRUD operations for routes, are supported. 
This feature is primarily intended for debugging routes created through Terraform.

Once the feature reaches General Availability, we plan to introduce support
for creating routing tables and attaching them to networks directly via the
CLI. Until then, we recommend users continue managing routing tables and 
attachments through the Terraform provider.

```
stackit beta routing-table [flags]
```

### Options

```
  -h, --help   Help for "stackit beta routing-table"
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

* [stackit beta](./stackit_beta.md)	 - Contains beta STACKIT CLI commands
* [stackit beta routing-table describe](./stackit_beta_routing-table_describe.md)	 - Describe a routing-table
* [stackit beta routing-table list](./stackit_beta_routing-table_list.md)	 - List all routing-tables
* [stackit beta routing-table route](./stackit_beta_routing-table_route.md)	 - Manage routes of a routing-table

