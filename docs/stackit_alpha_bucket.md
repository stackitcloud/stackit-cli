## stackit alpha bucket

Manages objects in STACKIT Object Storage buckets (data plane)

### Synopsis

Manages objects in STACKIT Object Storage buckets.

Buckets and objects can be addressed with s3 URIs (s3://BUCKET[/KEY]) or
STACKIT Object Storage URLs in virtual hosted style
(https://BUCKET.object.storage.REGION.onstackit.cloud[/KEY]) or path style
(https://object.storage.REGION.onstackit.cloud/BUCKET[/KEY]).

Credentials (access key ID / secret access key) are required and can be provided
with the --access-key-id and --secret-access-key flags, the STACKIT_OBJECT_STORAGE_ACCESS_KEY_ID
and STACKIT_OBJECT_STORAGE_SECRET_ACCESS_KEY environment variables, or the config file.
Create credentials with 'stackit object-storage credentials create'.

```
stackit alpha bucket [flags]
```

### Examples

```
  List all buckets
  $ stackit alpha bucket ls

  Upload a file to a bucket
  $ stackit alpha bucket cp local/file.txt s3://my-bucket/file.txt
```

### Options

```
  -h, --help[=true|false]   Help for "stackit alpha bucket"
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

* [stackit alpha](./stackit_alpha.md)	 - Contains alpha STACKIT CLI commands
* [stackit alpha bucket cp](./stackit_alpha_bucket_cp.md)	 - Copies objects between local and Object Storage
* [stackit alpha bucket ls](./stackit_alpha_bucket_ls.md)	 - Lists buckets and objects
* [stackit alpha bucket mb](./stackit_alpha_bucket_mb.md)	 - Creates a bucket
* [stackit alpha bucket mv](./stackit_alpha_bucket_mv.md)	 - Moves objects between local and Object Storage
* [stackit alpha bucket presign](./stackit_alpha_bucket_presign.md)	 - Generates a presigned URL for an object
* [stackit alpha bucket rb](./stackit_alpha_bucket_rb.md)	 - Deletes a bucket
* [stackit alpha bucket rm](./stackit_alpha_bucket_rm.md)	 - Deletes objects

