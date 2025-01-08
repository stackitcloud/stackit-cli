## stackit project create

Creates a STACKIT project

### Synopsis

Creates a STACKIT project.
You can associate a project with a STACKIT Network Area (SNA) by providing the ID of the SNA.
The STACKIT Network Area (SNA) allows projects within an organization to be connected to each other on a network level.
This makes it possible to connect various resources of the projects within an SNA and also simplifies the connection with on-prem environments (hybrid cloud).
The network type can no longer be changed after the project has been created. If you require a different network type, you must create a new project.


```
stackit project create [flags]
```

### Examples

```
  Create a STACKIT project
  $ stackit project create --parent-id xxxx --name my-project

  Create a STACKIT project with a set of labels
  $ stackit project create --parent-id xxxx --name my-project --label key=value --label foo=bar

  Create a STACKIT project with a network area
  $ stackit project create --parent-id xxxx --name my-project --network-area-id yyyy
```

### Options

```
  -h, --help                     Help for "stackit project create"
      --label stringToString     Labels are key-value string pairs which can be attached to a project. A label can be provided with the format key=value and the flag can be used multiple times to provide a list of labels (default [])
      --name string              Project name
      --network-area-id string   ID of a STACKIT Network Area (SNA) to associate with the project.
      --parent-id string         Parent resource identifier. Both container ID (user-friendly) and UUID are supported
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

* [stackit project](./stackit_project.md)	 - Manages projects

