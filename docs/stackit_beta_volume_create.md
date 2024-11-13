## stackit beta volume create

Creates a volume

### Synopsis

Creates a volume.

```
stackit beta volume create [flags]
```

### Examples

```
  Create a volume with availability zone "eu01-1"
  $ stackit beta volume create --availability-zone eu01-1

  Create a volume with name "volume-1"
  $ stackit beta volume create --availability-zone eu01-1 --name volume-1

  Create a volume with availability zone "eu01-1", performance class "storage_premium_perf1"
  $ stackit beta volume create --availability-zone eu01-1 --performance-class storage_premium_perf1
```

### Options

```
      --availability-zone string   Availability zone
      --description string         Volume description
  -h, --help                       Help for "stackit beta volume create"
      --label stringToString       Labels are key-value string pairs which can be attached to a volume. A label can be provided with the format key=value and the flag can be used multiple times to provide a list of labels (default [])
  -n, --name string                Volume name
      --performance-class string   Performance class
      --size int                   Volume size (GB). Either size or source ID and type flags must be given.
      --source-id string           ID of the source object of volume. Either size or source id and type flags must be given. Source ID and type must be provided together.
      --source-type string         Type of the source object of volume. Either size or source id and type flags must be given. Source ID and type must be provided together.
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

* [stackit beta volume](./stackit_beta_volume.md)	 - Provides functionality for Volume

