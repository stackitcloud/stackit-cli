## stackit beta scf organization role set

Assigns a role to a user of a STACKIT Cloud Foundry organization

### Synopsis

Assigns a role to a user of a STACKIT Cloud Foundry organization. The user is identified either by its GUID or by its username.

```
stackit beta scf organization role set ORGANIZATION_ID [flags]
```

### Examples

```
  Assign the role "organization_manager" to the user with GUID "yyy" in the STACKIT Cloud Foundry organization with ID "xxx"
  $ stackit beta scf organization role set xxx --role organization_manager --user-guid yyy

  Assign the role "organization_auditor" to the user with username "jane.doe" in the STACKIT Cloud Foundry organization with ID "xxx"
  $ stackit beta scf organization role set xxx --role organization_auditor --user-name jane.doe
```

### Options

```
  -h, --help               Help for "stackit beta scf organization role set"
      --role string        The role to assign to the user (one of: [organization_user, organization_auditor, organization_manager, organization_billing_manager])
      --user-guid string   The GUID of the user to assign the role to
      --user-name string   The username of the user to assign the role to
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

* [stackit beta scf organization role](./stackit_beta_scf_organization_role.md)	 - Manages roles of STACKIT Cloud Foundry organizations

