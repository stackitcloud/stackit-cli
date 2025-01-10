## stackit beta volume create

Creates a volume

### Synopsis

Creates a volume.

```
stackit beta volume create [flags]
```

### Examples

```
  Create a volume with availability zone "eu01-1" and size 64 GB
  $ stackit beta volume create --availability-zone eu01-1 --size 64

  Create a volume with availability zone "eu01-1", size 64 GB and labels
  $ stackit beta volume create --availability-zone eu01-1 --size 64 --labels key=value,foo=bar

  Create a volume with name "volume-1", from a source image with ID "xxx"
  $ stackit beta volume create --availability-zone eu01-1 --name volume-1 --source-id xxx --source-type image

  Create a volume with availability zone "eu01-1", performance class "storage_premium_perf1" and size 64 GB
  $ stackit beta volume create --availability-zone eu01-1 --performance-class storage_premium_perf1 --size 64
```

### Options

```
      --availability-zone string   Availability zone
      --description string         Volume description
  -h, --help                       Help for "stackit beta volume create"
      --labels stringToString      Labels are key-value string pairs which can be attached to a volume. E.g. '--labels key1=value1,key2=value2,...' (default [])
  -n, --name string                Volume name
      --performance-class string   Performance class
      --size int                   Volume size (GB). Either 'size' or the 'source-id' and 'source-type' flags must be given
      --source-id string           ID of the source object of volume. Either 'size' or the 'source-id' and 'source-type' flags must be given
      --source-type string         Type of the source object of volume. Either 'size' or the 'source-id' and 'source-type' flags must be given
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

* [stackit beta volume](./stackit_beta_volume.md)	 - Provides functionality for volumes

