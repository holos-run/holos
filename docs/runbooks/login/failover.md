## Overview

Use this run book to move the ZITADEL primary postgres database from one cluster to another.  At the end of the process https://zitadel.example.com service will be live on the cluster which was previously the standby at the start of the process.

Downtime: ~10 minutes if executed quickly.  ~60 minutes if executed slowly and deliberately.

## Initial State
1. https://login.ois.run is routed to core2.
2. `postgrescluster/zitadel` on core2 is primary.
3. `postgrescluster/zitadel` on core1 is standby.
4. A recent [[#Full Backup]] has been taken and is stored in S3.  
5. The standby cluster has been restored from the recent full backup.
## Process
1. [[#Edit Platform site config]]
2. [[#Render core2]]
3. [[#Demote core2]]  
4. [[#Verify core2 is standby]]
5. [[#Render core1]]
6. [[#Promote core1]]
7. [[#Verify core1 is primary]]
8. [[#Reconcile Zitadel]]
9. [[#Cut over NLB]]
10. [[#Test Login]]

## Edit Platform site config
Change the primary cluster from core2 to core1 in `platforms/platform.site.cue`.
```cue
#Platform: primaryCluster: name: "core1"
```
## Render core2
Configure core as a standby.
```
holos render --cluster-name=core2 ~/workspace/holos-run/holos/docs/examples/platforms/reference/clusters/accounts/...
```

Expected change:
```diff
diff --git a/deploy/clusters/core2/components/prod-iam-postgres/prod-iam-postgres.gen.yaml b/deploy/clusters/core2/components/prod-iam-postgres/prod-iam-postgres.gen.yaml
index c1aa8394..45daa650 100644
--- a/deploy/clusters/core2/components/prod-iam-postgres/prod-iam-postgres.gen.yaml
+++ b/deploy/clusters/core2/components/prod-iam-postgres/prod-iam-postgres.gen.yaml
@@ -85,24 +85,6 @@ spec:
     name: zitadel-repl-tls
   customTLSSecret:
     name: zitadel-primary-tls
-  dataSource:
-    pgbackrest:
-      configuration:
-        - secret:
-            name: pgo-s3-creds
-      global:
-        repo2-path: /pgbackrest/prod-iam-zitadel/zitadel/repo2
-        repo2-cipher-type: aes-256-cbc
-      options:
-        - --type=time
-        - --target="2024-03-11 17:10:00+00"
-      repo:
-        name: repo2
-        s3:
-          bucket: ois-zitadel-backups
-          endpoint: s3.dualstack.us-east-2.amazonaws.com
-          region: us-east-2
-      stanza: db
   image: registry.developers.crunchydata.com/crunchydata/crunchy-postgres:ubi8-16.2-0
   instances:
     - name: pgha1
@@ -139,7 +121,7 @@ spec:
       image: registry.developers.crunchydata.com/crunchydata/crunchy-pgbouncer:ubi8-1.21-3
       replicas: 2
   standby:
-    enabled: false
+    enabled: true
     repoName: repo2
   users:
     - databases
```
## Demote core2

```
10012  git add -p
10013  git commit -m 'Make core2 a standby zitadel cluster'
10014  git push origin HEAD
10015  flux reconcile source git flux-system
10016  flux get ks
```

## Verify core2 is standby
```
k logs $(kubectl get pods -o name -l postgres-operator.crunchydata.com/role=master) | tail -2
```
Expected output:
```
2024-03-14 21:58:26,205 INFO: no action. I am (zitadel-pgha1-smw7-0), the standby leader with the lock
```
## Render core1
```
holos render --cluster-name=core1 ~/workspace/holos-run/holos/docs/examples/platforms/reference/clusters/accounts/...
```
Expected change:
```diff
diff --git a/deploy/clusters/core1/components/prod-iam-postgres/prod-iam-postgres.gen.yaml b/deploy/clusters/core1/components/prod-iam-postgres/prod-iam-postgres.gen.yaml
index aa69814e..290c0865 100644
--- a/deploy/clusters/core1/components/prod-iam-postgres/prod-iam-postgres.gen.yaml
+++ b/deploy/clusters/core1/components/prod-iam-postgres/prod-iam-postgres.gen.yaml
@@ -85,6 +85,22 @@ spec:
     name: zitadel-repl-tls
   customTLSSecret:
     name: zitadel-primary-tls
+  dataSource:
+    pgbackrest:
+      configuration:
+        - secret:
+            name: pgo-s3-creds
+      global:
+        repo2-path: /pgbackrest/prod-iam-zitadel/zitadel/repo2
+        repo2-cipher-type: aes-256-cbc
+      options: []
+      repo:
+        name: repo2
+        s3:
+          bucket: ois-zitadel-backups
+          endpoint: s3.dualstack.us-east-2.amazonaws.com
+          region: us-east-2
+      stanza: db
   image: registry.developers.crunchydata.com/crunchydata/crunchy-postgres:ubi8-16.2-0
   instances:
     - name: pgha1
@@ -121,7 +137,7 @@ spec:
       image: registry.developers.crunchydata.com/crunchydata/crunchy-pgbouncer:ubi8-1.21-3
       replicas: 2
   standby:
-    enabled: true
+    enabled: false
     repoName: repo2
   users:
     - databases:
```

## Promote core1
```
flux reconcile source git flux-system
flux reconcile ks prod-iam-postgres
```
## Verify core1 is primary
```
k logs $(kubectl get pods -o name -l postgres-operator.crunchydata.com/role=master) | tail -1
```
Expected output:
```
2024-03-14 22:05:01,159 INFO: no action. I am (zitadel-pgha1-xl2b-0), the leader with the lock
```
## Reconcile Zitadel
Now that the database is back up, reconcile zitadel.
This is pretty brutal but it works:
```
❯ flux delete ks prod-iam-zitadel
Are you sure you want to delete this kustomization: y
► deleting kustomization prod-iam-zitadel in flux-system namespace
✔ kustomization deleted
```
```
❯ k apply --server-side=true -f deploy/clusters/core1/holos/components/prod-iam-zitadel-kustomization.gen.yaml
kustomization.kustomize.toolkit.fluxcd.io/prod-iam-zitadel serverside-applied
```
ZITADEL should come up READY 2/2
```
❯ k get pods
NAME                               READY   STATUS      RESTARTS      AGE
zitadel-76776f747f-94wft           2/2     Running     2 (30s ago)   34s
zitadel-76776f747f-fl24c           2/2     Running     1 (30s ago)   34s
zitadel-76776f747f-xx2w7           2/2     Running     2 (30s ago)   34s
zitadel-backup-49c7-h25cj          0/1     Completed   0             3m23s
zitadel-init-td2sh                 0/1     Completed   0             34s
zitadel-pgbouncer-d9f8cffc-2vwk5   2/2     Running     0             22m
zitadel-pgbouncer-d9f8cffc-rrrhd   2/2     Running     0             22m
zitadel-pgha1-4npq-0               4/4     Running     0             3m21s
zitadel-pgha1-xl2b-0               4/4     Running     0             3m43s
zitadel-repo-host-0                2/2     Running     0             22m
zitadel-setup-b7zkx                0/1     Completed   1             34s
zitadel-test-connection            0/1     Error       0             34s
```

## Cut over NLB
Finally, cut over the NLB config to move to core1.
This is specific to our metal clusters.
Apply with: `git push puppet +HEAD:production -o run=netlb`
```diff
diff --git a/site-modules/profile/templates/haproxy/haproxy.cfg.epp b/site-modules/profile/templates/haproxy/haproxy.cfg.epp
index b358a70..3ee40db 100644
--- a/site-modules/profile/templates/haproxy/haproxy.cfg.epp
+++ b/site-modules/profile/templates/haproxy/haproxy.cfg.epp
@@ -225,7 +225,7 @@ frontend sni-ingress from tcp-ingress
   use_backend core2-https if { req.ssl_sni -i web.holos.run }
   use_backend core2-https if { req.ssl_sni -i -m end .web.holos.run }
   # Identity provider is routed to one cluster at a time
-  use_backend core2-https if { req.ssl_sni -i login.ois.run }
+  use_backend core1-https if { req.ssl_sni -i login.ois.run }
   # Holos Dev Environments
   use_backend k1-https if { req.ssl_sni -i holos.wtf }
   use_backend k1-https if { req.ssl_sni -i -m end .holos.wtf }
@@ -358,7 +358,7 @@ frontend http-ingress from http-ingress
   bind 65.102.23.41:80
   bind 2602:41:6617:2ec0::80:1d:80 v6only
   # Zitadel is active on one cluster at a time
-  use_backend core2-http if { hdr(host) -i login.ois.run }
+  use_backend core1-http if { hdr(host) -i login.ois.run }
   # Vault is active on core1 or core2
   use_backend core2-http if { hdr(host) -i vault.core.ois.run }
   # Infrastructure Dex OIDC ID issuer is active on core1 or core2

```
## Test Login
```
kubectl oidc-login get-token \
  --oidc-extra-scope=openid \
  --oidc-extra-scope=email \
  --oidc-extra-scope=profile \
  --oidc-extra-scope=groups \
  --oidc-extra-scope=offline_access \
  --oidc-extra-scope=urn:zitadel:iam:org:domain:primary:openinfrastructure.co \
  --oidc-issuer-url=https://login.ois.run \
  --oidc-client-id=${CLIENT_ID:-257714027772314751@holos_platform} \
  --oidc-use-pkce \
  --force-refresh > ~/.kube/token.json
```

```
jq -r .status.token ~/.kube/token.json | cut -d. -f2 | base64 -d | jq
```

```json
{
  "amr": [
    "mfa",
    "user"
  ],
  "at_hash": "5nUq6kXWT8GheCNvDNdyow",
  "aud": [
    "257714027772314751@holos_platform",
    "257713952794870157"
  ],
  "auth_time": 1710172186,
  "azp": "257714027772314751@holos_platform",
  "client_id": "257714027772314751@holos_platform",
  "email": "jeff@openinfrastructure.co",
  "email_verified": true,
  "exp": 1710497614,
  "family_name": "McCune",
  "given_name": "Jeff",
  "groups": [
    "prod-cluster-admin"
  ],
  "iat": 1710454404,
  "iss": "https://login.ois.run",
  "locale": null,
  "name": "Jeff McCune",
  "nickname": "Jeff",
  "preferred_username": "jeff@openinfrastructure.co",
  "sub": "257712562366383231",
  "updated_at": 1710105084,
  "urn:zitadel:iam:org:domain:primary": "openinfrastructure.co",
  "urn:zitadel:iam:org:project:257713952794870157:roles": {
    "prod-cluster-admin": {
      "257712562366317695": "openinfrastructure.co"
    }
  },
  "urn:zitadel:iam:org:project:roles": {
    "prod-cluster-admin": {
      "257712562366317695": "openinfrastructure.co"
    }
  }
}
```
## Tasks
### Identify the leader
Same command for primary and standby clusters.
```
k get pods -o name -l postgres-operator.crunchydata.com/role=master
```
### Full Backup
On the primary cluster:
```
kubectl -n prod-iam-zitadel annotate postgrescluster zitadel postgres-operator.crunchydata.com/pgbackrest-backup="$(date)" --overwrite
```
Watch the progress:
```
k -n prod-iam-zitadel logs -l postgres-operator.crunchydata.com/pgbackrest-backup=manual -f
```
Expected output:
```
time="2024-03-14T18:36:16Z" level=info msg="crunchy-pgbackrest starts"
time="2024-03-14T18:36:16Z" level=info msg="debug flag set to false"
time="2024-03-14T18:36:16Z" level=info msg="backrest backup command requested"
time="2024-03-14T18:36:16Z" level=info msg="command to execute is [pgbackrest backup --stanza=db --repo=2 --type=full]"
time="2024-03-14T18:39:11Z" level=info msg="output=[]"
time="2024-03-14T18:39:11Z" level=info msg="stderr=[]"
time="2024-03-14T18:39:11Z" level=info msg="crunchy-pgbackrest ends"
```

## Troubleshooting
Full backup doesn't start
Check the controller logs.  If you see `Switchover failed` errors, try editing the number of postgrescluster replicas down to 1 then back up to 2.  Then try a backup again.
```
Switchover failed, details: 503, Switchover failed\n" version=5.5.1-0-amd64
```
