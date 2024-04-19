FROM registry.choria.io/choria/provisioner:latest

RUN curl -Lo nsc.zip https://github.com/nats-io/nsc/releases/download/v2.8.6/nsc-linux-amd64.zip &&\
  unzip nsc.zip && \
  mv nsc /usr/local/bin/nsc && \
  chmod 755 /usr/local/bin/nsc && \
  rm -f nsc.zip

# TODO: Add jwt executable
# TODO: Add helper executable

USER choria
ENV USER=choria

ENTRYPOINT ["/usr/sbin/choria-provisioner"]

# These two files are expected to be in the provisioner secret.
CMD ["--config=/etc/provisioner/provisioner.yaml", "--choria-config=/etc/provisioner/choria.cfg"]
