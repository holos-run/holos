package holos

import (
	cpv1 "pkg.crossplane.io/provider/v1"
	cpdrcv1beta1 "pkg.crossplane.io/deploymentruntimeconfig/v1beta1"
	cpfuncv1beta1 "pkg.crossplane.io/function/v1beta1"
	cpawspcv1beta1 "aws.upbound.io/providerconfig/v1beta1"
)

// Produce a helm chart build plan.
(#Helm & Chart).BuildPlan

// https://github.com/crossplane/crossplane/releases
let CrossplaneVersion = "1.16.0"

// https://github.com/crossplane-contrib/provider-upjet-aws/releases
let AWSProviderVersion = "v1.5.0"

let Chart = {
	Name:      "crossplane"
	Version:   CrossplaneVersion
	Namespace: "crossplane-system"

	Repo: name: "crossplane-stable"
	Repo: url:  "https://charts.crossplane.io/stable"

	Values: {
		podSecurityContextCrossplane: {
			runAsNonRoot: true
			seccompProfile: type: "RuntimeDefault"
		}
		securityContextCrossplane: capabilities: drop: ["ALL"]
		podSecurityContextRBACManager: {
			runAsNonRoot: true
			seccompProfile: type: "RuntimeDefault"
		}
		securityContextRBACManager: capabilities: drop: ["ALL"]
	}

	Resources: {
		// This DeploymentRuntimeConfig adds annotations to the service accounts spun up for the AWS providers.
		// https://docs.crossplane.io/latest/concepts/providers/#configuring-metadata-of-runtime-resources
		//
		// Adding this SA annotation causes the EKS Pod Identity Webhook to inject
		// environment variables for AWS authentication as well as mount the AWS token
		// file to the provider pod that uses the SA. For example:
		//
		//  Environment:
		//     AWS_STS_REGIONAL_ENDPOINTS:   regional
		//     AWS_DEFAULT_REGION:           us-east-1
		//     AWS_REGION:                   us-east-1
		//     AWS_ROLE_ARN:                 arn:aws:iam::271053619184:role/holos-crossplane
		//     AWS_WEB_IDENTITY_TOKEN_FILE:  /var/run/secrets/eks.amazonaws.com/serviceaccount/token
		//  Mounts:
		//       /var/run/secrets/eks.amazonaws.com/serviceaccount from aws-iam-token (ro)
		//
		// Docs: https://docs.crossplane.io/latest/concepts/providers/#runtime-configuration
		DeploymentRuntimeConfig: "aws-irsa": cpdrcv1beta1.#DeploymentRuntimeConfig & {
			metadata: name: "aws-irsa"
			spec: serviceAccountTemplate: metadata: annotations: "eks.amazonaws.com/role-arn": "arn:aws:iam::\(_Platform.Model.aws.accountNumber):role/holos-crossplane"
			spec: deploymentTemplate: spec: template: spec: containers: [
				{
					name: "package-runtime"
					args: ["--enable-external-secret-stores"]
				},
			]
		}

		// https://marketplace.upbound.io/providers/upbound/provider-family-aws
		Provider: "upbound-provider-family-aws": cpv1.#Provider & {
			metadata: name: "upbound-provider-family-aws"
			spec: package:  "xpkg.upbound.io/upbound/provider-family-aws:\(AWSProviderVersion)"
			// The provider-family-aws provider doesn't need the IRSA SA annotation.
			spec: runtimeConfigRef: name: "default"
		}
		// https://marketplace.upbound.io/providers/upbound/provider-aws-s3
		Provider: "provider-aws-s3": cpv1.#Provider & {
			metadata: name: "provider-aws-s3"
			spec: package:  "xpkg.upbound.io/upbound/provider-aws-s3:\(AWSProviderVersion)"
			spec: runtimeConfigRef: name: "aws-irsa"
		}
		// https://marketplace.upbound.io/providers/upbound/provider-aws-rds
		Provider: "provider-aws-rds": cpv1.#Provider & {
			metadata: name: "provider-aws-rds"
			spec: package:  "xpkg.upbound.io/upbound/provider-aws-rds:\(AWSProviderVersion)"
			spec: runtimeConfigRef: name: "aws-irsa"
		}

		// The patch-and-transform function is used in Compositions.
		// https://github.com/crossplane-contrib/function-patch-and-transform/releases
		Function: "function-patch-and-transform": cpfuncv1beta1.#Function & {
			metadata: name: "function-patch-and-transform"
			spec: package:  "xpkg.upbound.io/crossplane-contrib/function-patch-and-transform:v0.5.0"
		}

		// By setting the ProviderConfig's name to "default", all resources with an
		// apiVersion of "aws.upbound.io/v1beta1" will use this ProviderConfig unless
		// otherwise specified with a providerConfigRef.
		ProviderConfig: default: cpawspcv1beta1.#ProviderConfig & {
			metadata: name: "default"
			spec: credentials: source: "IRSA"
		}
	}
}
