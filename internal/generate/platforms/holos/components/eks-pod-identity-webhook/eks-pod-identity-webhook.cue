// This sets up the AWS EKS Pod Identity Webhook, which is used to inject AWS
// credentials into pods so services can use IRSA for AWS authentication.
//
// See: https://github.com/aws/amazon-eks-pod-identity-webhook
//
// There isn't an official Helm chart for the EKS Pod Identity Webhook, so we
// use https://github.com/jkroepke/helm-charts/tree/main/charts/amazon-eks-pod-identity-webhook
// See: https://github.com/aws/amazon-eks-pod-identity-webhook/issues/4

package holos

// https://github.com/jkroepke/helm-charts/tree/main/charts/amazon-eks-pod-identity-webhook
let ChartVersion = "2.1.3"

// https://github.com/aws/amazon-eks-pod-identity-webhook/releases
// https://registry.hub.docker.com/r/amazon/amazon-eks-pod-identity-webhook/tags
let AppVersion = "v0.5.4"

let Chart = {
	Name:      "amazon-eks-pod-identity-webhook"
	Version:   ChartVersion
	Namespace: "aws-pod-identity"

	Repo: name: "jkroepke"
	Repo: url:  "https://jkroepke.github.io/helm-charts"

	Values: {
		image: tag: AppVersion
		config: {
			tokenAudience:    "sts.amazonaws.com"
			defaultAwsRegion: _Platform.Model.aws.primaryRegion
			extraArgs: ["-v=4"] // verbosity of at least 4 is needed to see mutation events.
		}
		securityContext: {
			allowPrivilegeEscalation: false
			capabilities: drop: ["ALL"]
			runAsNonRoot: true
		}
		podSecurityContext: {
			// https://github.com/aws/amazon-eks-pod-identity-webhook/blob/master/README.md#usage-with-non-root-container-user
			fsGroup: 2000
			seccompProfile: type: "RuntimeDefault"
		}
	}
}

// Produce a helm chart build plan.
(#Helm & Chart).Output
