## stackit image create

Creates images

### Synopsis

Creates images.

```
stackit image create [flags]
```

### Examples

```
  Create an image with name 'my-new-image' from a raw disk image located in '/my/raw/image'
  $ stackit image create --name my-new-image --disk-format=raw --local-file-path=/my/raw/image

  Create an image with name 'my-new-image' from a qcow2 image read from '/my/qcow2/image' with labels describing its contents
  $ stackit image create --name my-new-image --disk-format=qcow2 --local-file-path=/my/qcow2/image --labels os=linux,distro=alpine,version=3.12

  Create an image with name 'my-new-image' from a raw disk image located in '/my/raw/image' with uefi disabled
  $ stackit image create --name my-new-image --disk-format=raw --local-file-path=/my/raw/image --uefi=false
```

### Options

```
      --architecture string      Sets the CPU architecture. By default x86 is used.
      --boot-menu                Enables the BIOS bootmenu.
      --cdrom-bus string         Sets CDROM bus controller type.
      --disk-bus string          Sets Disk bus controller type.
      --disk-format string       The disk format of the image. 
  -h, --help                     Help for "stackit image create"
      --labels stringToString    Labels are key-value string pairs which can be attached to a network-interface. E.g. '--labels key1=value1,key2=value2,...' (default [])
      --local-file-path string   The path to the local disk image file.
      --min-disk-size int        Size in Gigabyte.
      --min-ram int              Size in Megabyte.
      --name string              The name of the image.
      --nic-model string         Sets virtual nic model.
      --no-progress              Show no progress indicator for upload.
      --os string                Enables OS specific optimizations.
      --os-distro string         Operating System Distribution.
      --os-version string        Version of the OS.
      --protected                Protected VM.
      --rescue-bus string        Sets the device bus when the image is used as a rescue image.
      --rescue-device string     Sets the device when the image is used as a rescue image.
      --secure-boot              Enables Secure Boot.
      --uefi                     Enables UEFI boot. (default true)
      --video-model string       Sets Graphic device model.
      --virtio-scsi              Enables the use of VirtIO SCSI to provide block device access. By default instances use VirtIO Block.
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

* [stackit image](./stackit_image.md)	 - Manage server images

