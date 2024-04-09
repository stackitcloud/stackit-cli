## stackit organization member remove

Removes a member from an organization

### Synopsis

Removes a member from an organization.
A member is a combination of a subject (user, service account or client) and a role.
The subject is usually email address (for users) or name (for clients).

```
stackit organization member remove SUBJECT [flags]
```

### Examples

```
  Remove a member (user "someone@domain.com" with an "editor" role) from an organization
  $ stackit organization member remove someone@domain.com --organization-id xxx --role editor

  Remove a member (user "someone@domain.com" with a "reader" role) from an organization, along with all other roles of the subject that would stop the removal of the "reader" role
  $ stackit organization member remove someone@domain.com --organization-id xxx --role reader --force
```

### Options

```
      --force                    When true, removes other roles of the subject that would stop the removal of the requested role
  -h, --help                     Help for "stackit organization member remove"
      --organization-id string   The organization ID
      --role string              The role to be removed from the subject
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit organization member](./stackit_organization_member.md)	 - Provides functionality regarding organization members

