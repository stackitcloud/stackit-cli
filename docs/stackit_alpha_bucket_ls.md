## stackit alpha bucket ls

Lists buckets and objects

### Synopsis

Lists all buckets, or the contents of a bucket.

Without an argument, lists all buckets.
With a bucket URI (s3://BUCKET), lists the "directories" and objects at the bucket root.
With an object URI (s3://BUCKET/PREFIX), lists the "directories" and objects under the prefix.
With --recursive, lists all objects under the prefix flat, including all nested keys.

```
stackit alpha bucket ls [BUCKET_URI] [flags]
```

### Examples

```
  List all buckets
  $ stackit alpha bucket ls

  List the contents of a bucket
  $ stackit alpha bucket ls s3://my-bucket

  Recursively list all objects under a prefix
  $ stackit alpha bucket ls s3://my-bucket/path/to/ --recursive
```

### Options

```
      --access-key-id string       Object Storage access key ID. Can also be set with the object_storage_access_key_id config key or the STACKIT_OBJECT_STORAGE_ACCESS_KEY_ID environment variable.
      --endpoint string            Object Storage API endpoint. Defaults to https://object.storage.<region>.onstackit.cloud. Can also be set with the object_storage_data_plane_custom_endpoint config key or the STACKIT_OBJECT_STORAGE_DATA_PLANE_CUSTOM_ENDPOINT environment variable.
  -h, --help[=true|false]          Help for "stackit alpha bucket ls"
      --recursive[=true|false]     List all objects under the prefix recursively, including nested keys
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

