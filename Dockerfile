FROM registry.k8s.io/kubectl:v1.31.0 AS kubectl
# https://github.com/GoogleContainerTools/distroless
FROM golang:1.23 AS build

WORKDIR /go/src/app
COPY . .

RUN CGO_ENABLED=0 make install

# Install helm to /usr/local/bin/helm
# https://helm.sh/docs/intro/install/#from-script
RUN curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 \
  && chmod 700 get-helm-3 \
  && DESIRED_VERSION=v3.16.2 ./get-helm-3 \
  && rm -f get-helm-3

# distroless
FROM gcr.io/distroless/static-debian12 AS final
COPY --from=build /go/bin/holos / \
     --from=build /usr/local/bin/helm /usr/local/bin/helm \
     --from=kubectl /bin/kubectl /usr/local/bin/kubectl
