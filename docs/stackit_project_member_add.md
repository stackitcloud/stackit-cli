## stackit project member add

Adds a member to a project

### Synopsis

Adds a member to a project.
A member is a combination of a subject (user, service account or client) and a role.
The subject is usually email address for users or name in case of clients
For more details on the available roles, run:
  $ stackit project role list --project-id <PROJECT ID>

```
stackit project member add SUBJECT [flags]
```

### Examples

```
  Add a member to a project with the "reader" role
  $ stackit project member add someone@domain.com --project-id xxx --role reader
```

### Options

```
  -h, --help          Help for "stackit project member add"
      --role string   The role to add to the subject
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

* [stackit project member](./stackit_project_member.md)	 - Provides functionality for project members

