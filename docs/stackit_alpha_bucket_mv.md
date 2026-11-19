## stackit alpha bucket mv

Moves objects between local and Object Storage

### Synopsis

Moves a single object or a directory tree, then deletes the source.

One of SRC and DESTINATION must be a bucket URI (s3://BUCKET[/KEY]) or a STACKIT
Object Storage URL (https://BUCKET.object.storage.REGION.onstackit.cloud[/KEY])
and the other a local path.
Bucket-to-bucket moves use a server-side copy.
Without --recursive, a single object is moved.
With --recursive, all objects under the prefix (or all files under the local directory) are moved.

```
stackit alpha bucket mv SRC DESTINATION [flags]
```

### Examples

```
  Move a file to a bucket
  $ stackit alpha bucket mv local/file.txt s3://my-bucket/file.txt

  Move an object to a new key
  $ stackit alpha bucket mv s3://my-bucket/old/file.txt s3://my-bucket/new/file.txt

  Move a prefix to another bucket recursively
  $ stackit alpha bucket mv s3://src-bucket/prefix/ s3://dst-bucket/prefix/ --recursive
```

### Options

```
      --access-key-id string       Object Storage access key ID. Can also be set with the object_storage_access_key_id config key or the STACKIT_OBJECT_STORAGE_ACCESS_KEY_ID environment variable.
      --endpoint string            Object Storage API endpoint. Defaults to https://object.storage.<region>.onstackit.cloud. Can also be set with the object_storage_data_plane_custom_endpoint config key or the STACKIT_OBJECT_STORAGE_DATA_PLANE_CUSTOM_ENDPOINT environment variable.
  -h, --help[=true|false]          Help for "stackit alpha bucket mv"
      --recursive[=true|false]     Move a directory tree or all objects under a prefix
      --secret-access-key string   Object Storage secret access key. Can also be set with the object_storage_secret_access_key config key or the STACKIT_OBJECT_STORAGE_SECRET_ACCESS_KEY environment variable.
```

### Options inherited from parent commands

```
  -y, --assume-yes[=true|false]   If set, skips all confirmation prompts
      --async[=true|false]        If set, runs the command asynchronously
  -o, --output-format string      Output format, (one of: [json, pretty, none, yaml])
  -p, --project-id string         Project ID
      --region string             Target region for region-specific requests
      --verbosity string          Verbosity of the CLI, (one of: [debug, info, warning, error]) (default "info")
```

### SEE ALSO

* [stackit alpha bucket](./stackit_alpha_bucket.md)	 - Manages objects in STACKIT Object Storage buckets (data plane)

