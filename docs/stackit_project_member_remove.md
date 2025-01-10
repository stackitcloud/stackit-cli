## stackit project member remove

Removes a member from a project

### Synopsis

Removes a member from a project.
A member is a combination of a subject (user, service account or client) and a role.
The subject is usually email address for users or name in case of clients

```
stackit project member remove SUBJECT [flags]
```

### Examples

```
  Remove a member (user "someone@domain.com" with an "editor" role) from a project
  $ stackit project member remove someone@domain.com --project-id xxx --role editor

  Remove a member (user "someone@domain.com" with a "reader" role) from a project, along with all other roles of the subject that would stop the removal of the "reader" role
  $ stackit project member remove someone@domain.com --project-id xxx --role reader --force
```

### Options

```
      --force         When true, removes other roles of the subject that would stop the removal of the requested role
  -h, --help          Help for "stackit project member remove"
      --role string   The role to be removed from the subject
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

* [stackit project member](./stackit_project_member.md)	 - Manages project members

