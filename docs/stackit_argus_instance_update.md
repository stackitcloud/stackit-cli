## stackit argus instance update

Updates an Argus instance

### Synopsis

Updates an Argus instance.

```
stackit argus instance update INSTANCE_ID [flags]
```

### Examples

```
  Update the plan of an Argus instance with ID "xxx" by specifying the plan ID
  $ stackit argus instance update xxx --plan-id yyy

  Update the plan of an Argus instance with ID "xxx" by specifying the plan name
  $ stackit argus instance update xxx --plan-name yyy

  Update the name of an Argus instance with ID "xxx"
  $ stackit argus instance update xxx --name new-instance-name
```

### Options

```
  -h, --help               Help for "stackit argus instance update"
      --name string        Instance name
      --plan-id string     Plan ID
      --plan-name string   Plan name
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty"]
  -p, --project-id string      Project ID
```

### SEE ALSO

* [stackit argus instance](./stackit_argus_instance.md)	 - Provides functionality for Argus instances

