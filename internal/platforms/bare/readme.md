Bare Platform

| Folder     | Description                                                                                             |
| -          | -                                                                                                       |
| forms      | Contains Platform and Project form and model definitions                                                |
| platform   | Contains the Platform resource that defines how to render the configuration for all Platform Components |
| components | Contains BuildPlan resources which define how to render individual Platform Components                  |

## Forms

To populate the form, the platform must already be created in the Web UI:

```bash
platformId="018f36fb-e3ff-7f7f-a5d1-7ca2bf499e94"
cue export ./forms/platform/ --out json \
  | jq '{platform_id: "'$platformId'", fields: .spec.fields}' \
  | grpcurl -H "x-oidc-id-token: $(holos token)" -d @ \
    app.dev.k2.holos.run:443 \
    holos.v1alpha1.PlatformService.PutForm
```
