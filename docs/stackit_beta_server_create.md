## stackit beta server create

Creates a server

### Synopsis

Creates a server.

```
stackit beta server create [flags]
```

### Examples

```
  Create a server from an image with id xxx
  $ stackit beta server create --machine-type t1.1 --name server1 --image-id xxx

  Create a server with labels from an image with id xxx
  $ stackit beta server create --machine-type t1.1 --name server1 --image-id xxx --labels key=value,foo=bar

  Create a server with a boot volume
  $ stackit beta server create --machine-type t1.1 --name server1 --boot-volume-source-id xxx --boot-volume-source-type image --boot-volume-size 64

  Create a server with a boot volume from an existing volume
  $ stackit beta server create --machine-type t1.1 --name server1 --boot-volume-source-id xxx --boot-volume-source-type volume

  Create a server with a keypair
  $ stackit beta server create --machine-type t1.1 --name server1 --image-id xxx --keypair-name example

  Create a server with a network
  $ stackit beta server create --machine-type t1.1 --name server1 --image-id xxx --network-id yyy

  Create a server with a network interface
  $ stackit beta server create --machine-type t1.1 --name server1 --boot-volume-source-id xxx --boot-volume-source-type image --boot-volume-size 64 --network-interface-ids yyy

  Create a server with an attached volume
  $ stackit beta server create --machine-type t1.1 --name server1 --boot-volume-source-id xxx --boot-volume-source-type image --boot-volume-size 64 --volumes yyy

  Create a server with user data (cloud-init)
  $ stackit beta server create --machine-type t1.1 --name server1 --boot-volume-source-id xxx --boot-volume-source-type image --boot-volume-size 64 --user-data file("${path.module}/cloud-init.yaml")
```

### Options

```
      --affinity-group string                  The affinity group the server is assigned to
      --availability-zone string               The availability zone of the server
      --boot-volume-delete-on-termination      Delete the volume during the termination of the server. Defaults to false
      --boot-volume-performance-class string   Boot volume performance class
      --boot-volume-size source_type           The size of the boot volume in GB. Must be provided when source_type is `image`
      --boot-volume-source-id string           ID of the source object of boot volume. It can be either 'image-id' or 'volume-id'
      --boot-volume-source-type string         Type of the source object of boot volume. It can be either  'image' or 'volume'
  -h, --help                                   Help for "stackit beta server create"
      --image-id string                        The image ID to be used for an ephemeral disk on the server. Either image-id or boot volume is required
      --keypair-name string                    The name of the SSH keypair used during the server creation
      --labels stringToString                  Labels are key-value string pairs which can be attached to a server. E.g. '--labels key1=value1,key2=value2,...' (default [])
      --machine-type string                    Name of the type of the machine for the server. Possible values are documented in https://docs.stackit.cloud/stackit/en/virtual-machine-flavors-75137231.html
  -n, --name string                            Server name
      --network-id string                      ID of the network for the initial networking setup for the server creation
      --network-interface-ids strings          List of network interface IDs for the initial networking setup for the server creation
      --security-groups strings                The initial security groups for the server creation
      --service-account-emails strings         List of the service account mails
      --user-data string                       User data that is passed via cloud-init to the server
      --volumes strings                        The list of volumes attached to the server
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

* [stackit beta server](./stackit_beta_server.md)	 - Provides functionality for Server

