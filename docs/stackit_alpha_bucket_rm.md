## stackit alpha bucket rm

Deletes objects

### Synopsis

Deletes a single object, or all objects under a prefix with --recursive.

To delete an empty bucket, use 'stackit alpha bucket rb'.

```
stackit alpha bucket rm BUCKET_URI [flags]
```

### Examples

```
  Delete a single object
  $ stackit alpha bucket rm s3://my-bucket/file.txt

  Delete all objects under a prefix without prompting
  $ stackit alpha bucket rm s3://my-bucket/path/to/ --recursive --force
```

### Options

```
      --access-key-id string       Object Storage access key ID. Can also be set with the object_storage_access_key_id config key or the STACKIT_OBJECT_STORAGE_ACCESS_KEY_ID environment variable.
      --endpoint string            Object Storage API endpoint. Defaults to https://object.storage.<region>.onstackit.cloud. Can also be set with the object_storage_data_plane_custom_endpoint config key or the STACKIT_OBJECT_STORAGE_DATA_PLANE_CUSTOM_ENDPOINT environment variable.
      --force[=true|false]         Do not prompt for confirmation
  -h, --help[=true|false]          Help for "stackit alpha bucket rm"
      --recursive[=true|false]     Delete all objects under the prefix
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

