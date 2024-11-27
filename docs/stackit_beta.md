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
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit](./stackit.md)	 - Manage STACKIT resources using the command line
* [stackit beta network](./stackit_beta_network.md)	 - Provides functionality for Network
* [stackit beta network-area](./stackit_beta_network-area.md)	 - Provides functionality for STACKIT Network Area (SNA)
* [stackit beta network-interface](./stackit_beta_network-interface.md)	 - Provides functionality for Network Interface
* [stackit beta public-ip](./stackit_beta_public-ip.md)	 - Provides functionality for Public IP
* [stackit beta server](./stackit_beta_server.md)	 - Provides functionality for Server
* [stackit beta sqlserverflex](./stackit_beta_sqlserverflex.md)	 - Provides functionality for SQLServer Flex
* [stackit beta volume](./stackit_beta_volume.md)	 - Provides functionality for Volume

