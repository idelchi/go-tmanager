#!/bin/bash

folder=~/.local/bin
folder=.bin-linux-arm7

for file in ${folder}/*; do
  if [[ -x ${file} && ! -d ${file} ]]; then
    # Skip "traefik"
    if [[ ${file} == *traefik ]]; then
      continue
    fi

    echo "Running ${file}"
    "${file}" --version 2>/dev/null
  fi
done
