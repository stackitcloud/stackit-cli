## stackit beta

Contains beta STACKIT CLI commands

### Synopsis

Contains beta STACKIT CLI commands.
The commands under this group are still in a beta state, and functionality may be incomplete or have breaking changes.

```
stackit beta [flags]
```

### Examples

```
  See the currently available beta commands
  $ stackit beta --help

  Execute a beta command
  $ stackit beta MY_COMMAND
```

### Options

```
  -h, --help   Help for "stackit beta"
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
* [stackit beta alb](./stackit_beta_alb.md)	 - Manages application loadbalancers
* [stackit beta intake](./stackit_beta_intake.md)	 - Provides functionality for intake
* [stackit beta kms](./stackit_beta_kms.md)	 - Provides functionality for KMS
* [stackit beta logs](./stackit_beta_logs.md)	 - Provides functionality for Logs
* [stackit beta sfs](./stackit_beta_sfs.md)	 - Provides functionality for SFS (stackit file storage)
* [stackit beta sqlserverflex](./stackit_beta_sqlserverflex.md)	 - Provides functionality for SQLServer Flex

