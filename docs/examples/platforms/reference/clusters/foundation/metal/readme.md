# Metal Clusters

This cluster type is overlaid onto other cluster types to add services necessary outside of a cloud like GKE or EKS.  Ceph for PersistenVolumeClaim support on a Talos Proxmox cluster is the primary use case.

## Test Script

Test ceph is working with:

```bash
apply -n default -f-<<EOF
heredoc> apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: test
spec:
  accessModes:
    - ReadWriteOnce
  volumeMode: Filesystem
  resources:
    requests:
      storage: 1Gi
EOF
```
