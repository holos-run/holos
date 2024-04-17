Initialize machine room provisioning credentials

Setup Notes:

The holos server flag `--provisioner-seed` must match the issuer.seed value.
To get the correct value to configure for holos server:

    holos get secret choria-issuer --print-key=issuer.seed --namespace $NAMESPACE
