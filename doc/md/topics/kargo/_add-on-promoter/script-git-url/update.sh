#! /bin/bash
set -euo pipefail
[[ -z "${HOLOS_UPDATE_SCRIPTS:-}" ]] && exit 0
cat > "$1"
