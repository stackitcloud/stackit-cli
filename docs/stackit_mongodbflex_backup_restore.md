## stackit mongodbflex backup restore

Restores a MongoDB Flex instance from a backup

### Synopsis

Restores a MongoDB Flex instance from a backup of an instance or clones a MongoDB Flex instance from a point-in-time snapshot.
The backup is specified by a backup id and the point-in-time snapshot is specified by a timestamp.
The instance to apply the backup to can be specified, otherwise it will be the same as the backup.

```
stackit mongodbflex backup restore [flags]
```

### Examples

```
  Restores a MongoDB Flex instance with id "yyy" using backup with id "zzz"
  $ stackit mongodbflex backup restore --instance-id yyy --backup-id zzz

  Clone a MongoDB Flex instance with id "yyy" via point-in-time restore to timestamp "zzz"
  $ stackit mongodbflex backup restore --instance-id yyy --timestamp zzz

  Restores a MongoDB Flex instance with id "yyy" using backup from instance with id "zzz" with backup id "aaa"
  $ stackit mongodbflex backup restore --instance-id zzz --backup-instance-id yyy --backup-id aaa
```

### Options

```
      --backup-id string            Backup id
      --backup-instance-id string   Instance id of the target instance to restore the backup to
  -h, --help                        Help for "stackit mongodbflex backup restore"
      --instance-id string          Instance id
      --timestamp string            Timestamp of the snapshot to clone the instance from
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit mongodbflex backup](./stackit_mongodbflex_backup.md)	 - Provides functionality for MongoDB Flex instance backups

