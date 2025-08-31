FROM registry.k8s.io/kubectl:v1.33.4 AS kubectl
# https://github.com/GoogleContainerTools/distroless
FROM golang:1.24 AS build

WORKDIR /go/src/app
COPY . .

RUN CGO_ENABLED=0 make install
RUN CGO_ENABLED=0 go install sigs.k8s.io/kustomize/kustomize/v5

# Install helm to /usr/local/bin/helm
# https://helm.sh/docs/intro/install/#from-script
# https://holos.run/docs/v1alpha5/tutorial/setup/#dependencies
RUN curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 \
  && chmod 700 get_helm.sh \
  && DESIRED_VERSION=v3.16.2 ./get_helm.sh \
  && rm -f get_helm.sh

COPY --from=kubectl /bin/kubectl /usr/local/bin/

# Use debian slim instead of distroless to get package management.
FROM public.ecr.aws/docker/library/debian:13-slim AS final
COPY --from=build \
     /go/bin/holos \
     /go/bin/kustomize \
     /usr/local/bin/kubectl \
     /usr/local/bin/helm \
     /bin/

# Extra packages
# git - https://github.com/holos-run/holos/issues/440
RUN apt update && \
    apt install -y --no-install-recommends git && \
    apt clean && \
    rm -rf /var/lib/apt/lists/*

# Usage: docker run -v $(pwd):/app --workdir /app --rm -it quay.io/holos-run/holos holos render platform
CMD ["/bin/holos"]
