## stackit project describe

Shows details of a STACKIT project

### Synopsis

Shows details of a STACKIT project.

```
stackit project describe [flags]
```

### Examples

```
  Get the details of the configured STACKIT project
  $ stackit project describe

  Get the details of a STACKIT project by explicitly providing the project ID
  $ stackit project describe --project-id xxx

  Get the details of the configured STACKIT project, including details of the parent resources
  $ stackit project describe --include-parents
```

### Options

```
  -h, --help              Help for "stackit project describe"
      --include-parents   When true, the details of the parent resources will be included in the output
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit project](./stackit_project.md)	 - Provides functionality regarding projects

