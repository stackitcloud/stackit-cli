## stackit beta vpn connection create

Creates a VPN connection

### Synopsis

Creates a VPN connection.

```
stackit beta vpn connection create [flags]
```

### Examples

```
  Create a VPN connection
  $ stackit beta vpn connection create --gateway-id xxx --display-name my-connection --tunnel1-remote-address 1.2.3.4 --tunnel2-remote-address 5.6.7.8
```

### Options

```
      --display-name string                            Required: A user friendly name for the connection.
      --enabled                                        Enable the connection (default true)
      --gateway-id string                              Required: Gateway ID
  -h, --help                                           Help for "stackit beta vpn connection create"
      --labels stringToString                          Map of custom labels. Key and values must be a string with max 63 chars, start/end with alphanumeric. The key of a label follows the same rules as the LabelValue except that it cannot be empty. (example: foo=bar) (default [])
      --local-subnets strings                          Defaults to 0.0.0.0/0 for Route-based VPN configurations. Mandatory for Policy-based.
      --remote-subnets strings                         Defaults to 0.0.0.0/0 for Route-based VPN configurations. Mandatory for Policy-based.
      --static-routes strings                          Use this for route-based VPN.
      --tunnel1-bgp-remote-asn int                     Required: Tunnel 1 BGP Remote ASN.
                                                       ASN for private use (reserved by IANA), both 16Bit and 32Bit ranges are valid (RFC 6996).
      --tunnel1-peering-local-address string           Tunnel 1 Peering Local Address.
                                                       The peering object defines the point-to-point IP configuration for the Tunnel Interface. These addresses serve as next-hop identifiers and are used for BGP peering sessions and can be used in Static Route-Based connectivity.
      --tunnel1-peering-remote-address string          Tunnel 1 Peering Remote Address
      --tunnel1-phase1-dh-groups strings               Tunnel 1 Phase 1 DH Groups.
                                                       The Diffie-Hellman Group. Required, except if AEAD algorithms are selected. (possible values: [modp1024, modp2048, ecp256, ecp384, modp2048s256]) (default [])
      --tunnel1-phase1-encryption-algorithms strings   Required: Tunnel 1 Phase 1 Encryption Algorithms (possible values: [aes256, aes128gcm16, aes256gcm16]) (default [])
      --tunnel1-phase1-integrity-algorithms strings    Required: Tunnel 1 Phase 1 Integrity Algorithms (possible values: [sha1, sha2_256, sha2_384]) (default [])
      --tunnel1-phase1-rekey-time int                  Tunnel 1 Phase 1 Rekey Time.
                                                       Time to schedule a IKE re-keying (in seconds).
      --tunnel1-phase2-dh-groups strings               Tunnel 1 Phase 2 DH Groups (possible values: [modp1024, modp2048, ecp256, ecp384, modp2048s256]) (default [])
      --tunnel1-phase2-dpd-action string               Tunnel 1 Phase 2 DPD Action.
                                                       Action to perform for this CHILD_SA on DPD timeout. "clear": Closes the CHILD_SA and does not take further action. "restart": immediately tries to re-negotiate the CILD_SA under a fresh IKE_SA. (possible values: [clear, restart])
      --tunnel1-phase2-encryption-algorithms strings   Required: Tunnel 1 Phase 2 Encryption Algorithms (possible values: [aes256, aes128gcm16, aes256gcm16]) (default [])
      --tunnel1-phase2-integrity-algorithms strings    Required: Tunnel 1 Phase 2 Integrity Algorithms (possible values: [sha1, sha2_256, sha2_384]) (default [])
      --tunnel1-phase2-rekey-time int                  Tunnel 1 Phase 2 Rekey Time.
                                                       Time to schedule a Child SA re-keying (in seconds).
      --tunnel1-phase2-start-action string             Tunnel 1 Phase 2 Start Action.
                                                       Action to perform after loading the connection configuration. "none": The connection will be loaded but needs to be manually initiated. "start": initiates the connection actively. (possible values: [none, start])
      --tunnel1-pre-shared-key string                  Required: Tunnel 1 Pre Shared Key.
                                                       A Pre-Shared Key for authentication. Required in create-requests, optional in update-requests and omitted in every response.
      --tunnel1-remote-address string                  Tunnel 1 Remote Address
      --tunnel2-bgp-remote-asn int                     Tunnel 2 BGP Remote ASN
      --tunnel2-peering-local-address string           Tunnel 2 Peering Local Address.
                                                       The peering object defines the point-to-point IP configuration for the Tunnel Interface. These addresses serve as next-hop identifiers and are used for BGP peering sessions and can be used in Static Route-Based connectivity.
      --tunnel2-peering-remote-address string          Tunnel 2 Peering Remote Address
      --tunnel2-phase1-dh-groups strings               Tunnel 2 Phase 1 DH Groups
                                                       The Diffie-Hellman Group. Required, except if AEAD algorithms are selected. (possible values: [modp1024, modp2048, ecp256, ecp384, modp2048s256]) (default [])
      --tunnel2-phase1-encryption-algorithms strings   Required: Tunnel 2 Phase 1 Encryption Algorithms (possible values: [aes256, aes128gcm16, aes256gcm16]) (default [])
      --tunnel2-phase1-integrity-algorithms strings    Required: Tunnel 2 Phase 1 Integrity Algorithms (possible values: [sha1, sha2_256, sha2_384]) (default [])
      --tunnel2-phase1-rekey-time int                  Tunnel 2 Phase 1 Rekey Time.
                                                       Time to schedule a IKE re-keying (in seconds).
      --tunnel2-phase2-dh-groups strings               Tunnel 2 Phase 2 DH Groups (possible values: [modp1024, modp2048, ecp256, ecp384, modp2048s256]) (default [])
      --tunnel2-phase2-dpd-action string               Tunnel 2 Phase 2 DPD Action.
                                                       Action to perform for this CHILD_SA on DPD timeout. "clear": Closes the CHILD_SA and does not take further action. "restart": immediately tries to re-negotiate the CILD_SA under a fresh IKE_SA. (possible values: [clear, restart])
      --tunnel2-phase2-encryption-algorithms strings   Required: Tunnel 2 Phase 2 Encryption Algorithms (possible values: [aes256, aes128gcm16, aes256gcm16]) (default [])
      --tunnel2-phase2-integrity-algorithms strings    Required: Tunnel 2 Phase 2 Integrity Algorithms (possible values: [sha1, sha2_256, sha2_384]) (default [])
      --tunnel2-phase2-rekey-time int                  Tunnel 2 Phase 2 Rekey Time.
                                                       Time to schedule a Child SA re-keying (in seconds).
      --tunnel2-phase2-start-action string             Tunnel 2 Phase 2 Start Action.
                                                       Default: "start"
                                                       Enum: "none" "start"
                                                       Action to perform after loading the connection configuration. "none": The connection will be loaded but needs to be manually initiated. "start": initiates the connection actively. (possible values: [none, start])
      --tunnel2-pre-shared-key string                  Required: Tunnel 2 Pre Shared Key.
                                                       A Pre-Shared Key for authentication. Required in create-requests, optional in update-requests and omitted in every response.
      --tunnel2-remote-address string                  Tunnel 2 Remote Address
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

* [stackit beta vpn connection](./stackit_beta_vpn_connection.md)	 - Provides functionality for VPN connections

