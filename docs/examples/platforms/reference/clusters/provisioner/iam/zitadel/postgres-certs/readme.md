# Database Certs

This component issues postgres certificates from the provisioner cluster using certmanager.

The purpose is to define customTLSSecret and customReplicationTLSSecret to provide certs that allow the standby to authenticate to the primary. For this type of standby, you must use custom TLS.

Refer to the PGO [Streaming Standby](https://access.crunchydata.com/documentation/postgres-operator/latest/tutorials/backups-disaster-recovery/disaster-recovery#streaming-standby) tutorial.
