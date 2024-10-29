#!/bin/sh
set -e

DISABLE_SSL=""
for arg in "$@"; do
  if [ "${arg}" = "-k" ]; then
    DISABLE_SSL=yes
    break
  fi
done

curl ${DISABLE_SSL:+-k} -sSL https://raw.githubusercontent.com/idelchi/scripts/refs/heads/dev/install.sh | INSTALLER_TOOL="godyl" sh -s -- "$@"
