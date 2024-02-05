## stackit project member list

List members of a project

### Synopsis

List members of a project.

```
stackit project member list [flags]
```

### Examples

```
  List all members of a project
  $ stackit project role list --project-id xxx

  List all members of a project, sorted by role
  $ stackit project role list --project-id xxx --sort-by role

  List up to 10 members of a project
  $ stackit project role list --project-id xxx --limit 10
```

### Options

```
  -h, --help             Help for "stackit project member list"
      --limit int        Maximum number of entries to list
      --sort-by string   Sort entries by a specific field, one of ["subject" "role"] (default "subject")
      --subject string   Filter by subject (Identifier of user, service account or client. Usually email address in case of users or name in case of clients)
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit project member](./stackit_project_member.md)	 - Provides functionality regarding project members

