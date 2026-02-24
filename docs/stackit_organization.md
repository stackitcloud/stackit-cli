## stackit organization

Manages organizations

### Synopsis

Manages organizations.
An active STACKIT organization is the root element of the resource hierarchy and a prerequisite to use any STACKIT Cloud Resource / Service.

```
stackit organization [flags]
```

### Options

```
  -h, --help   Help for "stackit organization"
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

* [stackit](./stackit.md)	 - Manage STACKIT resources using the command line
* [stackit organization describe](./stackit_organization_describe.md)	 - Show an organization
* [stackit organization list](./stackit_organization_list.md)	 - Lists all organizations
* [stackit organization member](./stackit_organization_member.md)	 - Manages organization members
* [stackit organization role](./stackit_organization_role.md)	 - Manages organization roles

