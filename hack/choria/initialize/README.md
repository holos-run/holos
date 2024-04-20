Initialize machine room provisioning credentials

When you want the holos controller to provision while operating in the current
working directory, run:

 1. `init-choria-provisioner-creds` to populate secrets in the Holos
    Provisioner Cluster (not to be confused with the Choria Provisioner).
 2. `make-provisioning-jwt` to issue a `provisioning.jwt` file for `holos
    controller` to use.
 3. `holos controller --config=agent.cfg` to read `provisioning.jwt` and write
    the provisioned config file and credentials to the current directory.

Expect the controller to provision.

Setup Notes:

The holos server flag `--provisioner-seed` must match the issuer.seed value.
To get the correct value to configure for holos server:

    holos get secret choria-issuer --print-key=issuer.seed --namespace $NAMESPACE
