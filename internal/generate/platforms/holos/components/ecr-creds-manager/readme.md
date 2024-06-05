# ECR Credentials Manager

This component manages a `ecr-creds-refresher` `CronJob` in the `holos-system` `Namespace` of the Management Cluster.  This job authenticates to AWS using workload identity.  Refer to [Use workload identity with AWS](https://cloud.google.com/kubernetes-engine/multi-cloud/docs/aws/how-to/use-workload-identity-aws) for information on how to configure AWS to accepts kubernetes service account tokens from the GKE Management Cluster.

Refer also to [Pod Identity](https://github.com/aws/amazon-eks-pod-identity-webhook?tab=readme-ov-file#eks-walkthrough)

> [!NOTE]
> Both documents refer to EKS, but the process is the same and works on any kubernetes cluster.

Example [trust policy][trust-policy]:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "Federated": "arn:aws:iam::637423192589:oidc-provider/container.googleapis.com/v1/projects/holos-ops/locations/us-central1/clusters/management"
            },
            "Action": "sts:AssumeRoleWithWebIdentity",
            "Condition": {
                "StringEquals": {
                    "container.googleapis.com/v1/projects/holos-ops/locations/us-central1/clusters/management:aud": "sts.amazonaws.com",
                    "container.googleapis.com/v1/projects/holos-ops/locations/us-central1/clusters/management:sub": "system:serviceaccount:holos-system:ecr-creds-manager"
                }
            }
        }
    ]
}
```

[trust-policy]: https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_condition-logic-multiple-context-keys-or-values.html
