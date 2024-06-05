FROM quay.io/holos-run/debian:bullseye AS final
USER root
WORKDIR /app
ADD bin bin
RUN chown -R app: /app
# Kubernetes requires the user to be numeric
USER 8192
ENTRYPOINT bin/holos server
