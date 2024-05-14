# PG Bouncer

Every few days ZITADEL fails.  The problem seems to be related to pgbouncer not
being able to resolve DNS.  Restarting the pgbouncer pod fixes the issue.

See [How to load-balance queries between several servers?](https://www.pgbouncer.org/faq.html#how-to-load-balance-queries-between-several-servers)

> [!NOTE]
> DNS round-robin. Use several IPs behind one DNS name. PgBouncer does not look up DNS each time a new connection is launched. Instead, it caches all IPs and does round-robin internally. Note: if there are more than 8 IPs behind one name, the DNS backend must support the EDNS0 protocol. See README for details.

## Workaround

```sh
# Get the tls based creds to bypass oidc
(cd ~/.kube && holos get secret core2-kubeconfig-admin --print-key kubeconfig.admin > core2.admin)
export KUBECONFIG=$HOME/.kube/core2.admin
# Restart pgbouncer
kubectl -n prod-iam rollout restart deployment zitadel-pgbouncer
```

## Symptom logs

```sh
kubectl -n prod-iam logs -c pgbouncer -l postgres-operator.crunchydata.com/role=pgbouncer
```

```txt
2024-05-08 17:56:11.424 UTC [7] LOG S-0x559b03f90ff0: zitadel/zitadel@10.110.109.110:5432 SSL established: TLSv1.3/TLS_AES_256_GCM_SHA384/ECDH=prime256v1
2024-05-08 17:56:11.429 UTC [7] LOG S-0x559b03f92820: zitadel/zitadel@10.110.109.110:5432 new connection to server (from 10.244.5.38:53658)
2024-05-08 17:56:11.435 UTC [7] LOG S-0x559b03f92820: zitadel/zitadel@10.110.109.110:5432 SSL established: TLSv1.3/TLS_AES_256_GCM_SHA384/ECDH=prime256v1
2024-05-08 17:56:11.476 UTC [7] LOG C-0x559b03f7a610: zitadel/zitadel@10.244.2.89:34932 closing because: client close request (age=440s)
2024-05-08 17:56:19.708 UTC [7] LOG stats: 15 xacts/s, 42 queries/s, 0 client parses/s, 0 server parses/s, 0 binds/s, in 6159 B/s, out 6124 B/s, xact 3930 us, query 869 us, wait 490 us
[msg] Nameserver 10.96.0.10:53 is back up
2024-05-08 17:57:09.366 UTC [7] LOG C-0x559b03f7a610: zitadel/zitadel@10.244.3.187:58674 login attempt: db=zitadel user=zitadel tls=TLSv1.3/TLS_AES_256_GCM_SHA384
2024-05-08 17:57:09.391 UTC [7] LOG C-0x559b03f7a610: zitadel/zitadel@10.244.3.187:58674 closing because: client close request (age=0s)
2024-05-08 17:57:19.709 UTC [7] LOG stats: 9 xacts/s, 24 queries/s, 0 client parses/s, 0 server parses/s, 0 binds/s, in 2870 B/s, out 3018 B/s, xact 4147 us, query 958 us, wait 23 us
2024-05-08 17:58:19.708 UTC [7] LOG stats: 12 xacts/s, 32 queries/s, 0 client parses/s, 0 server parses/s, 0 binds/s, in 3861 B/s, out 3533 B/s, xact 3843 us, query 853 us, wait 0 us
2024-05-08 17:56:11.411 UTC [8] LOG S-0x55a894e36650: zitadel/_crunchypgbouncer@10.110.109.110:5432 new connection to server (from 10.244.3.227:58984)
2024-05-08 17:56:11.411 UTC [8] LOG S-0x55a894e37920: zitadel/zitadel@10.110.109.110:5432 new connection to server (from 10.244.3.227:58992)
2024-05-08 17:56:11.418 UTC [8] LOG S-0x55a894e37920: zitadel/zitadel@10.110.109.110:5432 SSL established: TLSv1.3/TLS_AES_256_GCM_SHA384/ECDH=prime256v1
2024-05-08 17:56:11.420 UTC [8] LOG S-0x55a894e36650: zitadel/_crunchypgbouncer@10.110.109.110:5432 SSL established: TLSv1.3/TLS_AES_256_GCM_SHA384/ECDH=prime256v1
2024-05-08 17:56:11.438 UTC [8] LOG S-0x55a894e35b90: zitadel/zitadel@10.110.109.110:5432 new connection to server (from 10.244.3.227:59004)
2024-05-08 17:56:11.445 UTC [8] LOG S-0x55a894e35b90: zitadel/zitadel@10.110.109.110:5432 SSL established: TLSv1.3/TLS_AES_256_GCM_SHA384/ECDH=prime256v1
2024-05-08 17:56:17.148 UTC [8] LOG stats: 9 xacts/s, 27 queries/s, 0 client parses/s, 0 server parses/s, 0 binds/s, in 3236 B/s, out 2826 B/s, xact 5224 us, query 910 us, wait 1182 us
[msg] Nameserver 10.96.0.10:53 is back up
2024-05-08 17:57:17.145 UTC [8] LOG stats: 10 xacts/s, 31 queries/s, 0 client parses/s, 0 server parses/s, 0 binds/s, in 4342 B/s, out 4305 B/s, xact 4536 us, query 776 us, wait 0 us
2024-05-08 17:58:17.149 UTC [8] LOG stats: 5 xacts/s, 15 queries/s, 0 client parses/s, 0 server parses/s, 0 binds/s, in 1641 B/s, out 1582 B/s, xact 7819 us, query 1426 us, wait 0 us
```

## Relevant Configuration

`/etc/pgbouncer/pgbouncer.ini` is empty.

```
bash-4.4$ cat /etc/pgbouncer/~postgres-operator.ini
# Generated by postgres-operator. DO NOT EDIT.
# Your changes will not be saved.

[pgbouncer]
%include /etc/pgbouncer/pgbouncer.ini

[pgbouncer]
auth_file = /etc/pgbouncer/~postgres-operator/users.txt
auth_query = SELECT username, password from pgbouncer.get_auth($1)
auth_user = _crunchypgbouncer
client_tls_ca_file = /etc/pgbouncer/~postgres-operator/frontend-ca.crt
client_tls_cert_file = /etc/pgbouncer/~postgres-operator/frontend-tls.crt
client_tls_key_file = /etc/pgbouncer/~postgres-operator/frontend-tls.key
client_tls_sslmode = require
conffile = /etc/pgbouncer/~postgres-operator.ini
ignore_startup_parameters = extra_float_digits
listen_addr = *
listen_port = 5432
server_tls_ca_file = /etc/pgbouncer/~postgres-operator/backend-ca.crt
server_tls_sslmode = verify-full
unix_socket_dir =

[databases]
* = host=zitadel-primary port=5432
```

### [host](https://www.pgbouncer.org/config.html#host)

> Host name or IP address to connect to. Host names are resolved at connection time, the result is cached per dns_max_ttl parameter. When a host name’s resolution changes, existing server connections are automatically closed when they are released (according to the pooling mode), and new server connections immediately use the new resolution. If DNS returns several results, they are used in a round-robin manner.

### dns_max_ttl

[dns_max_ttl](https://www.pgbouncer.org/config.html#dns_max_ttl)

How long DNS lookups can be cached. The actual DNS TTL is ignored.

Default: 15.0 (seconds)