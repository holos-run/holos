FROM 271053619184.dkr.ecr.us-east-2.amazonaws.com/holos-run/container-images/debian:bullseye AS final
USER root
WORKDIR /app
ADD bin bin
RUN chown -R app: /app
USER app
ENTRYPOINT bin/holos server
