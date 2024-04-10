#! /bin/bash

set -euo pipefail

PARENT="$(cd $(dirname "$0") && pwd)"

# If necessary
if [[ -s "${PARENT}/aws-login.last" ]]; then
  last="$(<"${PARENT}/aws-login.last")"
  now="$(date +%s)"
  if [[ $(( now - last )) -lt 28800 ]]; then
    echo "creds are still valid" >&2
    exit 0
  fi
fi

aws sso logout
aws sso login
aws ecr get-login-password --region us-east-2 | docker login --username AWS --password-stdin "${AWS_ACCOUNT}.dkr.ecr.${AWS_DEFAULT_REGION}.amazonaws.com"
# Touch a file so tilt docker_build can watch it as a dep
date +%s > "${PARENT}/aws-login.last"
