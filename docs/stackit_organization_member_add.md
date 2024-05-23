## stackit organization member add

Adds a member to an organization

### Synopsis

Adds a member to an organization.
A member is a combination of a subject (user, service account or client) and a role.
The subject is usually email address for users or name in case of clients
For more details on the available roles, run:
  $ stackit organization role list --organization-id <RESOURCE ID>

```
stackit organization member add SUBJECT [flags]
```

### Examples

```
  Add a member to an organization with the "reader" role
  $ stackit organization member add someone@domain.com --organization-id xxx --role reader
```

### Options

```
  -h, --help                     Help for "stackit organization member add"
      --organization-id string   The organization ID
      --role string              The role to add to the subject
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit organization member](./stackit_organization_member.md)	 - Provides functionality for organization members

