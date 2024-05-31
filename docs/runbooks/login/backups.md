# PostgresCluster Backups

This document describes how the S3 bucket for `PostgresCluster` backups is configured.  These buckets are configured both for ZITADEL and for Holos 
Server and are applicable to any service in Holos that stores data in a pgo `PostgresCluster` resource.

## Create the Bucket
Name: `holos-zitadel-backups` for `zitadel`
Name: `holos-server-backups` for `holos server`
> [!NOTE]
> The settings below match the default settings recommended by AWS.

Object Ownership: `ACLs disabled` (recommended) Checked.
Block Public Access settings for this bucket: **`Block all public access`** Checked.
Bucket Versioning: `Disable`
Default encryption: `Server-side encryption with Amazon S3 managed keys (SSE-S3)`
Bucket Key: `Enable`
Object Lock: `Disable`

## Create an IAM Policy
Create one IAM Policy for each bucket to grant full access to the bucket.  Replace the resource with each bucket name.
Name: `holos-zitadel-backups` for `zitadel`
Name: `holos-server-backups` for `holos server`
Description: `Read and write access to a specific bucket for pgrest operating within a pgo PostgresCluster.`

Policy JSON:
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "s3:GetBucketLocation",
                "s3:ListAllMyBuckets"
            ],
            "Resource": "arn:aws:s3:::*"
        },
        {
            "Effect": "Allow",
            "Action": "s3:*",
            "Resource": [
                "arn:aws:s3:::holos-zitadel-backups",
                "arn:aws:s3:::holos-zitadel-backups/*"
            ]
        }
    ]
}
```
## Create an IAM Group
Create an IAM Group to attach the policy granting access to the bucket.
Name: `holos-zitadel-backups` for `zitadel`
Attach permission policies: `holos-zitadel-backups`

Name: `holos-server-backups` for `holos server`
Attach permission policies: `holos-server-backups`
## Create the IAM User
Create an IAM User entity for each PostgresCluster.  Do not provide user access to the AWS Management Console.
Name: `holos-zitadel-backups` for `zitadel`
Group: `holos-zitadel-backups`

Name: `holos-server-backups` for `holos server`
Group: `holos-server-backups`

## Create an Access Key
Create an access key for `pgbackrest` associated with the `PostgresCluster`.

Description: 
> Used by pgbackrest associated with the PostgresCluster resource.  Refer to the PostgresCluster resource pgbackrest.cofiguration.secret.name for the stored location of the access key.  Synced from the Management Cluster using an ExternalSecret.
## Create the Secret
Create a `Secret` in the holos management cluster usable by pgbackrest.  This is a secret with a single key, `s3.conf` with the following format:
```
[global]
repo2-cipher-pass=
repo2-s3-key=
repo2-s3-key-secret=
repo3-cipher-pass=
repo3-s3-key=
repo3-s3-key-secret=
```
> [!NOTE]
> Use the same values for repo2 and repo3.  The purpose is to make space for migrating if need be in the future.

Generate the cipher pass using.  This password is used to encrypt all backups using client side before the backup is written to the bucket.
```
tr -dc A-Za-z0-9 </dev/urandom | head -c 64
```

Store the secret into the management cluster:
```
holos create secret --namespace zitadel holos-zitadel-backups \
  --append-hash=false --from-file .
```

```
holos create secret --namespace holos holos-server-backups \
  --append-hash=false --from-file .
```
