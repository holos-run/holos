Initialize machine room provisioning credentials

Wanting the holos controller to provision operating in the current working
directory, run:

 1. init-choria-provisioner-creds
 2. make-provisioner-jwt
 3. `holos controller --config=agent`

Expect the controller to provision.

Setup Notes:

The holos server flag `--provisioner-seed` must match the issuer.seed value.
To get the correct value to configure for holos server:

    holos get secret choria-issuer --print-key=issuer.seed --namespace $NAMESPACE
