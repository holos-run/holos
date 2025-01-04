#! /bin/bash
set -euo pipefail
[[ -s "$1" ]] && [[ -z "${HOLOS_UPDATE_SCRIPTS:-}" ]] && exit 0
cat > "$1"
