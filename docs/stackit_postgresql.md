## stackit postgresql

Provides functionality for PostgreSQL

### Synopsis

Provides functionality for PostgreSQL

### Examples

```
$ stackit postgresql instance create --project-id xxx --name my-instance --plan-name plan-name --version version
$ stackit postgresql instance list --project-id xxx
$ stackit postgresql credential create --project-id xxx --instance-id xxx
$ stackit postgresql credential describe --project-id xxx --instance-id xxx --credentials-id xxx
```

### Options

```
  -h, --help   help for postgresql
```

### Options inherited from parent commands

```
      --project-id string   Project ID
```

### SEE ALSO

* [stackit](./stackit.md)	 - The root command of the STACKIT CLI
* [stackit postgresql credential](./stackit_postgresql_credential.md)	 - Provides functionality for PostgreSQL credentials
* [stackit postgresql instance](./stackit_postgresql_instance.md)	 - Provides functionality for PostgreSQL instance
* [stackit postgresql offerings](./stackit_postgresql_offerings.md)	 - Provides information regarding the PostgreSQL service offerings

