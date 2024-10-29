#!/bin/sh
set -e

INSTALL_DIR="./bin"
DISABLE_SSL=""

# Usage function
usage() {
    cat <<EOF
Usage: ${0} [OPTIONS]
Installs Kubernetes-related tools using godyl.

This script will install:
- helm
- kubectl (with alias 'kc')
- k9s
- kubectx
- kubens
- task

Output directory can be controlled with the '-o' flag. Defaults to './bin'.

Example:

  curl -sSL https://raw.githubusercontent.com/idelchi/godyl/refs/heads/dev/scripts/k8s.sh | sh -s

EOF
    printf "Options:\n"

    printf "  -o DIR\tOutput directory for installed tools (default: ./bin)\n"
    printf "  -k    \tDisable SSL verification\n"

    # curl ${DISABLE_SSL:+-k} -sSL https://raw.githubusercontent.com/idelchi/scripts/refs/heads/dev/install.sh | INSTALLER_TOOL="godyl" sh -s -- -p

    exit 1
}

# Parse arguments
parse_args() {
    while getopts ":o:h" opt; do
        case "${opt}" in
            o) INSTALL_DIR="${OPTARG}" ;;
            k) DISABLE_SSL=yes ;;
            h) usage ;;
        esac
    done
}

# Create and handle temporary directory
setup_temp_dir() {
    if [ -z "${TEMP_DIR}" ]; then
        TEMP_DIR=$(mktemp -d)
        debug "Created temporary directory: ${TEMP_DIR}"
    else
        mkdir -p "${TEMP_DIR}"
        debug "Using specified temporary directory: ${TEMP_DIR}"
    fi

    # Set trap to clean up temporary directory
    trap 'rm -rf "${TEMP_DIR}"' EXIT
}

# Install godyl and tools
install_tools() {
    tmp=$(mktemp -d)
    # trap 'rm -rf "${tmp}"' EXIT

    curl ${DISABLE_SSL:+-k} -sSL "https://raw.githubusercontent.com/idelchi/scripts/refs/heads/dev/install.sh" | INSTALLER_TOOL=godyl sh -s -- -d "${tmp}" ${DISABLE_SSL:+-k}
    printf "godyl installed to ${tmp}\n"

    printf "Installing tools to ${INSTALL_DIR}\n"

    # Install tools using godyl
    "${tmp}/godyl" ${DISABLE_SSL:+-k} --output="${INSTALL_DIR}" - <<YAML
- name: helm/helm
  path: https://get.helm.sh/helm-{{ .Version }}-{{ .OS }}-{{ .ARCH }}.tar.gz
- name: kubernetes/kubernetes
  exe: kubectl
  path: https://dl.k8s.io/{{ .Version }}/bin/{{ .OS }}/{{ .ARCH }}/kubectl{{ .EXTENSION }}
  aliases: kc
- derailed/k9s
- name: ahmetb/kubectx
- name: ahmetb/kubectx
  exe: kubens
- name: go-task/task
YAML

    rm -rf ${tmp}
    printf "All tools installed successfully to ${INSTALL_DIR}\n"
}

need_cmd() {
    if ! command -v "${1}" >/dev/null 2>&1; then
        printf "Required command '${1}' not found"
        exit 1
    fi
}

main() {
    parse_args "$@"

    # Check for required commands
    need_cmd curl

    # Install tools
    install_tools
}

main "$@"
