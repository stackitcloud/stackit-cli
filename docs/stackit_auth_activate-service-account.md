## stackit auth activate-service-account

Activate service account authentication

### Synopsis

Activate authentication using service account credentials.
For more details on how to configure your service account, check the Authentication section on our documentation (LINK HERE README)

```
stackit auth activate-service-account [flags]
```

### Examples

```
$ stackit auth activate-service-account --service-account-key-path path/to/service_account_key.json --private-key-path path/to/private_key.pem
```

### Options

```
  -h, --help                              help for activate-service-account
      --jwks-custom-endpoint string       Custom endpoint for the jwks API, which is used to get the json web key sets (jwks) to validate tokens when the service-account authentication is activated
      --private-key-path string           RSA private key path
      --service-account-key-path string   Service account key path
      --service-account-token string      Service account long-lived access token
      --token-custom-endpoint string      Custom endpoint for the token API, which is used to request access tokens when the service-account authentication is activated
```

### Options inherited from parent commands

```
      --project-id string   Project ID
```

### SEE ALSO

* [stackit auth](./stackit_auth.md)	 - Provides authentication functionality

