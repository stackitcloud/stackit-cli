## stackit alpha bucket presign

Generates a presigned URL for an object

### Synopsis

Generates a presigned URL that grants time-limited read access to an object
without exposing the Object Storage credentials.

```
stackit alpha bucket presign BUCKET_URI [flags]
```

### Examples

```
  Generate a presigned URL valid for 1 hour (default)
  $ stackit alpha bucket presign s3://my-bucket/file.txt

  Generate a presigned URL valid for 15 minutes
  $ stackit alpha bucket presign s3://my-bucket/file.txt --expires 15m
```

### Options

```
      --access-key-id string       Object Storage access key ID. Can also be set with the object_storage_access_key_id config key or the STACKIT_OBJECT_STORAGE_ACCESS_KEY_ID environment variable.
      --endpoint string            Object Storage API endpoint. Defaults to https://object.storage.<region>.onstackit.cloud. Can also be set with the object_storage_data_plane_custom_endpoint config key or the STACKIT_OBJECT_STORAGE_DATA_PLANE_CUSTOM_ENDPOINT environment variable.
      --expires duration           Expiration of the presigned URL (e.g. 15m, 1h, 72h). Maximum: 7 days (default 1h0m0s)
  -h, --help[=true|false]          Help for "stackit alpha bucket presign"
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

