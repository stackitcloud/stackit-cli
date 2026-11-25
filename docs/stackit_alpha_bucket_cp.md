## stackit alpha bucket cp

Copies objects between local and Object Storage

### Synopsis

Copies a single object or a directory tree between local storage and Object Storage.

One of SRC and DESTINATION must be a bucket URI (s3://BUCKET[/KEY]) or a STACKIT
Object Storage URL (https://BUCKET.object.storage.REGION.onstackit.cloud[/KEY])
and the other a local path.
Without --recursive, a single object is transferred.
With --recursive, all objects under the prefix (download) or all files under the
local directory (upload) are transferred.

```
stackit alpha bucket cp SRC DESTINATION [flags]
```

### Examples

```
  Upload a file to a bucket
  $ stackit alpha bucket cp local/file.txt s3://my-bucket/file.txt

  Download a single object
  $ stackit alpha bucket cp s3://my-bucket/file.txt local/file.txt

  Upload a directory recursively
  $ stackit alpha bucket cp ./local/dir s3://my-bucket/prefix/ --recursive

  Download a prefix recursively
  $ stackit alpha bucket cp s3://my-bucket/prefix/ ./local/dir --recursive
```

### Options

```
      --access-key-id string       Object Storage access key ID. Can also be set with the object_storage_access_key_id config key or the STACKIT_OBJECT_STORAGE_ACCESS_KEY_ID environment variable.
      --endpoint string            Object Storage API endpoint. Defaults to https://object.storage.<region>.onstackit.cloud. Can also be set with the object_storage_data_plane_custom_endpoint config key or the STACKIT_OBJECT_STORAGE_DATA_PLANE_CUSTOM_ENDPOINT environment variable.
  -h, --help[=true|false]          Help for "stackit alpha bucket cp"
      --recursive[=true|false]     Transfer a directory tree or all objects under a prefix
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

