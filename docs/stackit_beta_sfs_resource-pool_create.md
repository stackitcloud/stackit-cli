## stackit beta sfs resource-pool create

Creates a SFS resource pool

### Synopsis

Creates a SFS resource pool.

The available performance class values can be obtained by running:
 $ stackit beta sfs performance-class list

```
stackit beta sfs resource-pool create [flags]
```

### Examples

```
  Create a SFS resource pool
  $ stackit beta sfs resource-pool create --availability-zone eu01-m --ip-acl 10.88.135.144/28 --performance-class Standard --size 500 --name resource-pool-01

  Create a SFS resource pool, allow only a single IP which can mount the resource pool
  $ stackit beta sfs resource-pool create --availability-zone eu01-m --ip-acl 250.81.87.224/32 --performance-class Standard --size 500 --name resource-pool-01

  Create a SFS resource pool, allow multiple IP ACL which can mount the resource pool
  $ stackit beta sfs resource-pool create --availability-zone eu01-m --ip-acl "10.88.135.144/28,250.81.87.224/32" --performance-class Standard --size 500 --name resource-pool-01

  Create a SFS resource pool with visible snapshots
  $ stackit beta sfs resource-pool create --availability-zone eu01-m --ip-acl 10.88.135.144/28 --performance-class Standard --size 500 --name resource-pool-01 --snapshots-visible
```

### Options

```
      --availability-zone string   Availability zone
  -h, --help                       Help for "stackit beta sfs resource-pool create"
      --ip-acl strings             List of network addresses in the form <address/prefix>, e.g. 192.168.10.0/24 that can mount the resource pool readonly (default [])
      --name string                Name
      --performance-class string   Performance class
      --size int                   Size of the pool in Gigabytes
      --snapshots-visible          Set snapshots visible and accessible to users
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

* [stackit beta sfs resource-pool](./stackit_beta_sfs_resource-pool.md)	 - Provides functionality for SFS resource pools

