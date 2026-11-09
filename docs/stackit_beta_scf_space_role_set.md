## stackit beta scf space role set

Assigns a role to a user of a STACKIT Cloud Foundry space

### Synopsis

Assigns a role to a user of a STACKIT Cloud Foundry space. The user is identified either by its GUID or by its username.

```
stackit beta scf space role set SPACE_ID [flags]
```

### Examples

```
  Assign the role "space_developer" to the user with GUID "zzz" in the STACKIT Cloud Foundry space with ID "yyy" of the organization with ID "xxx"
  $ stackit beta scf space role set yyy --organization-id xxx --role space_developer --user-guid zzz

  Assign the role "space_auditor" to the user with username "jane.doe" in the STACKIT Cloud Foundry space with ID "yyy" of the organization with ID "xxx"
  $ stackit beta scf space role set yyy --organization-id xxx --role space_auditor --user-name jane.doe
```

### Options

```
  -h, --help[=true|false]        Help for "stackit beta scf space role set"
      --organization-id string   The UUID of the organization
      --role string              The role to assign to the user (one of: [space_auditor, space_developer, space_manager, space_supporter])
      --user-guid string         The GUID of the user to assign the role to
      --user-name string         The username of the user to assign the role to
```

### Options inherited from parent commands

```
  -y, --assume-yes[=true|false]   If set, skips all confirmation prompts
      --async[=true|false]        If set, runs the command asynchronously
  -o, --output-format string      Output format, (one of: [json, pretty, none, yaml])
  -p, --project-id string         Project ID
      --region string             Target region for region-specific requests
      --verbosity string          Verbosity of the CLI, (one of: [debug, info, warning, error]) (default "info")
```

### SEE ALSO

* [stackit beta scf space role](./stackit_beta_scf_space_role.md)	 - Manages roles of STACKIT Cloud Foundry spaces

