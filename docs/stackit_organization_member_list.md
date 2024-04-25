## stackit organization member list

Lists members of an organization

### Synopsis

Lists members of an organization

```
stackit organization member list [flags]
```

### Examples

```
  List all members of an organization
  $ stackit organization member list --organization-id xxx

  List all members of an organization in JSON format
  $ stackit organization member list --organization-id xxx --output-format json

  List up to 10 members of an organization
  $ stackit organization member list --organization-id xxx --limit 10
```

### Options

```
  -h, --help                     Help for "stackit organization member list"
      --limit int                Maximum number of entries to list
      --organization-id string   The organization ID
      --sort-by string           Sort entries by a specific field, one of ["subject" "role"] (default "subject")
      --subject string           Filter by subject (Identifier of user, service account or client. Usually email address in case of users or name in case of clients)
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit organization member](./stackit_organization_member.md)	 - Provides functionality regarding organization members

