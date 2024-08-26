## KubeStart - Kick start your software development platform

<img width="50%"
align="right"
style="display: block; margin: 40px auto;"
src="https://openinfrastructure.co/blog/2016/02/27/logo/logorectangle.png">

Building and maintaining a software development platform is a complex and time
consuming endeavour.  Organizations often dedicate a team of 3-4 who need 6-12
months to build the platform.

KubeStart is a collection of tooling and API specifications to reduce the
complexity and speed up the process of building a modern, cloud native software
development platform.

- **Accelerate new projects** - Reduce time to market and operational
  complexity by integrating software distributed with Helm, Kustomize, or any
  other tool that produces YAML into the KubeStart toolchain.
- **Modernize existing projects** - Incrementally incorporate your existing
  platform services with the KubeStart toolchain.
- **Unified configuration model** - Increase safety and reduce the risk of
  config changes with CUE.
- **First class Helm and Kustomize support** - Leverage and reuse your existing
  investment in existing configuration tools such as Helm and Kustomize.

## Quick Installation

```console
go install github.com/kube-start/kubestart@latest
```

## Docs and Support

The documentation for developing and using KubeStart is available at: https://kubestart.org

For discussion and support, [open a
discussion](https://github.com/kube-start/kubestart/discussions/new/choose).

## License

KubeStart is licensed under Apache 2.0 as found in the [LICENSE file](LICENSE).
