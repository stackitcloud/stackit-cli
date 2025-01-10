## stackit mongodbflex backup restore

Restores a MongoDB Flex instance from a backup

### Synopsis

Restores a MongoDB Flex instance from a backup of an instance or clones a MongoDB Flex instance from a point-in-time backup.
The backup can be specified by a backup ID or a timestamp.
You can specify the instance to which the backup will be applied. If not specified, the backup will be applied to the same instance from which it was taken.

```
stackit mongodbflex backup restore [flags]
```

### Examples

```
  Restore a MongoDB Flex instance with ID "yyy" using backup with ID "zzz"
  $ stackit mongodbflex backup restore --instance-id yyy --backup-id zzz

  Clone a MongoDB Flex instance with ID "yyy" via point-in-time restore to timestamp "2024-05-14T14:31:48Z"
  $ stackit mongodbflex backup restore --instance-id yyy --timestamp 2024-05-14T14:31:48Z

  Restore a MongoDB Flex instance with ID "yyy", using backup from instance with ID "zzz" with backup ID "xxx"
  $ stackit mongodbflex backup restore --instance-id zzz --backup-instance-id yyy --backup-id xxx
```

### Options

```
      --backup-id string            Backup ID
      --backup-instance-id string   Instance ID of the target instance to restore the backup to
  -h, --help                        Help for "stackit mongodbflex backup restore"
      --instance-id string          Instance ID
      --timestamp string            Timestamp to restore the instance to, in a date-time with the RFC3339 layout format, e.g. 2024-01-01T00:00:00Z
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

* [stackit mongodbflex backup](./stackit_mongodbflex_backup.md)	 - Provides functionality for MongoDB Flex instance backups

