# Postgres Full Backup

Suppose you delete all objects in the S3 bucket hosting all postgres backups.  You want to take a full backup ASAP of an existing PostgreSQL database.

The normal method of annotating the `postgrescluster` resource will not work because the job will error:

```
❯ kubectl annotate postgrescluster zitadel postgres-operator.crunchydata.com/pgbackrest-backup="$(date)" --overwrite
postgrescluster.postgres-operator.crunchydata.com/zitadel annotated
```

Backup fails:

```
❯ k get pods
NAME                                    READY   STATUS      RESTARTS      AGE
zitadel-backup-hk7w-76bfk               0/1     Error       0             65s
zitadel-backup-hk7w-d55v6               0/1     Error       0             44s
zitadel-backup-hk7w-l9dwm               0/1     Error       0             76s
zitadel-backup-hk7w-zcg69               0/1     Error       0             3s
zitadel-pgbouncer-d9f8cffc-nx8lq        2/2     Running     0             49m
zitadel-pgbouncer-d9f8cffc-s7g7x        2/2     Running     0             49m
zitadel-pgha1-2xv2-0                    5/5     Running     0             48m
zitadel-pgha1-78f4-0                    5/5     Running     0             49m
zitadel-repo-host-0                     2/2     Running     0             49m
```

Error is: `FileMissingError: unable to open missing file '/pgbackrest/prod-iam/zitadel/repo2/backup/db/backup.info.copy' for read`

```
time="2024-04-08T00:02:11Z" level=info msg="crunchy-pgbackrest starts"
time="2024-04-08T00:02:11Z" level=info msg="debug flag set to false"
time="2024-04-08T00:02:12Z" level=info msg="backrest backup command requested"
time="2024-04-08T00:02:12Z" level=info msg="command to execute is [pgbackrest backup --stanza=db --repo=2 --type=full]"
time="2024-04-08T00:02:12Z" level=info msg="output=[]"
time="2024-04-08T00:02:12Z" level=info msg="stderr=[ERROR: [055]: unable to load info file '/pgbackrest/prod-iam/zitadel/repo2/backup/db/backup.info' or '/pgbackrest/prod-iam/zitadel/repo2/backup/db/backup.info.copy':\n       FileMissingError: unable to open missing file '/pgbackrest/prod-iam/zitadel/repo2/backup/db/backup.info' for read\n       FileMissingError: unable to open missing file '/pgbackrest/prod-iam/zitadel/repo2/backup/db/backup.info.copy' for read\n       HINT: backup.info cannot be opened and is required to perform a backup.\n       HINT: has a stanza-create been performed?\n]"
time="2024-04-08T00:02:12Z" level=fatal msg="command terminated with exit code 55"
```

## Fix Process

We need to edit the postgrescluster.  We're going to have the controller re-initialize the backup repository from scratch by removing it and re-adding it.

First, suspend flux:

```
flux suspend ks prod-iam-zitadel prod-iam-postgres
```

Save the config to two files:

```
kubectl get postgresclusters.postgres-operator.crunchydata.com zitadel -o yaml > orig.yaml
cp orig.yaml new.yaml
```

Remove the follow fields and re-apply the cluster.  This will leave the cluster running and available while the controller reconciles the repo configuration:

```diff
--- orig.yaml	2024-04-07 17:08:26.834715820 -0700
+++ new.yaml	2024-04-07 17:08:57.418546067 -0700
@@ -4,6 +4,4 @@
   annotations:
     holos.run/description: ""
-    postgres-operator.crunchydata.com/pgbackrest-backup: Sun 07 Apr 2024 05:01:35
-      PM PDT
   creationTimestamp: "2024-04-07T23:10:44Z"
   finalizers:
@@ -26,12 +24,5 @@
         repo1-retention-full: "1"
         repo2-cipher-type: aes-256-cbc
-        repo2-path: /pgbackrest/prod-iam/zitadel/repo2
-        repo2-retention-full: "14"
-        repo2-retention-full-type: time
       image: registry.developers.crunchydata.com/crunchydata/crunchy-pgbackrest:ubi8-2.49-0
-      manual:
-        options:
-        - --type=full
-        repoName: repo2
       repos:
       - name: repo1
@@ -43,12 +34,4 @@
               requests:
                 storage: 4Gi
-      - name: repo2
-        s3:
-          bucket: ois-zitadel-backups
-          endpoint: s3.dualstack.us-east-2.amazonaws.com
-          region: us-east-2
-        schedules:
-          differential: 0 1 * * 1-6
-          full: 0 1 * * 0
       restore:
         enabled: true
```

Apply the config and wait for the controller to reconcile:

```
k apply --server-side=true -f new.yaml --force-conflicts
```

Check for reconciliation:

```
kubectl -n postgres-operator logs -l app.kubernetes.io/name=pgo | tail -1
```

```
time="2024-04-08T00:10:03Z" level=debug msg="reconciled cluster" controller=postgrescluster controllerGroup=postgres-operator.crunchydata.com controllerKind=PostgresCluster name=zitadel namespace=prod-iam postgresCluster=prod-iam/zitadel reconcileID=cc8c8eb7-9787-4504-8ecd-a04ec84fbc0b version=5.5.1-0-amd64
```

Re-add the repo host configuration

```
grep -v 'resourceVersion:' orig.yaml | k apply --server-side=true --force-conflicts -f-
```

```
postgrescluster.postgres-operator.crunchydata.com/zitadel serverside-applied
```

The full backup should be running and writing to S3 now:

```
kubectl logs -l postgres-operator.crunchydata.com/pgbackrest-backup=manual
```

```
time="2024-04-08T00:12:54Z" level=info msg="crunchy-pgbackrest starts"
time="2024-04-08T00:12:54Z" level=info msg="debug flag set to false"
time="2024-04-08T00:12:54Z" level=info msg="backrest backup command requested"
time="2024-04-08T00:12:54Z" level=info msg="command to execute is [pgbackrest backup --stanza=db --repo=2 --type=full]"
time="2024-04-08T00:16:02Z" level=info msg="output=[]"
time="2024-04-08T00:16:02Z" level=info msg="stderr=[]"
time="2024-04-08T00:16:02Z" level=info msg="crunchy-pgbackrest ends"
```

Finally, resume flux:

```
flux resume ks prod-iam-postgres prod-iam-zitadel
```

```
► resuming kustomization prod-iam-postgres in flux-system namespace
✔ kustomization resumed
► resuming kustomization prod-iam-zitadel in flux-system namespace
✔ kustomization resumed
```
