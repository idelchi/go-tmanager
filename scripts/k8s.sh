#!/bin/sh
set -e

# Allow setting via environment variables, will be overridden by flags
INSTALL_DIR=${GODYL_INSTALL_DIR:-"${HOME}/.local/bin"}
TEMP_DIR=${GODYL_TEMP_DIR}
DEBUG=${GODYL_DEBUG:-0}
DRY_RUN=${GODYL_DRY_RUN:-0}
DISABLE_SSL=${GODYL_DISABLE_SSL}
GODYL_VERSION=${GODYL_GODYL_VERSION:-"v0.2-beta"}

# Output formatting
format_message() {
    local color="${1}"
    local message="${2}"
    local prefix="${3}"

    # Only use colors if output is a terminal
    if [ -t 1 ]; then
        case "${color}" in
            red)    printf '\033[0;31m%s\033[0m\n' "${prefix}${message}" >&2 ;;
            yellow) printf '\033[0;33m%s\033[0m\n' "${prefix}${message}" >&2 ;;
            green)  printf '\033[0;32m%s\033[0m\n' "${prefix}${message}" ;;
            *)      printf '%s\n' "${prefix}${message}" ;;
        esac
    else
        printf '%s\n' "${prefix}${message}"
    fi
}

debug() {
    if [ "${DEBUG}" -eq 1 ]; then
        format_message "yellow" "$*" "DEBUG: "
    fi
}

warning() {
    format_message "red" "$*" "Warning: "
}

info() {
    format_message "" "$*"
}

success() {
    format_message "green" "$*"
}

# Check if a command exists
need_cmd() {
    if ! command -v "${1}" >/dev/null 2>&1; then
        warning "Required command '${1}' not found"
        exit 1
    fi
    debug "Found required command: ${1}"
}

# Usage function
usage() {
    cat <<EOF
Usage: ${0} [OPTIONS]
Installs Kubernetes-related tools using godyl.

Flags and environment variables:
    Flag  Env                Default              Description
    -----------------------------------------------------------------
    -d    GODYL_INSTALL_DIR   "${HOME}/.local/bin" Installation directory
    -t    GODYL_TEMP_DIR      <auto>              Temporary directory
    -v    GODYL_GODYL_VERSION "v0.2-beta"         Godyl version to use
    -x    GODYL_DEBUG                             Enable debug output
    -n    GODYL_DRY_RUN                           Dry run mode
    -k    GODYL_DISABLE_SSL                       Disable SSL verification, when set to non-empty value
    -h                                           Show this help message

Flags take precedence over environment variables when both are set.

Example:
    ${0} -d /usr/local/bin -x

This script will install:
- helm
- kubectl (with alias 'kc')
- k9s
- kubectx
- kubens
- task

Example:

  curl -sSL https://raw.githubusercontent.com/idelchi/godyl/refs/heads/dev/scripts/k8s.sh | sh -s -- -k

EOF
    exit 1
}

# Parse arguments
parse_args() {
    while getopts ":d:t:v:xnkh" opt; do
        case "${opt}" in
            d) INSTALL_DIR="${OPTARG}" ;;
            t) TEMP_DIR="${OPTARG}" ;;
            v) GODYL_VERSION="${OPTARG}" ;;
            x) DEBUG=1 ;;
            n) DRY_RUN=1 ;;
            k) DISABLE_SSL=1 ;;
            h) usage ;;
            :) warning "Option -${OPTARG} requires an argument"; usage ;;
            *) warning "Invalid option: -${OPTARG}"; usage ;;
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
    debug "Installing godyl version ${GODYL_VERSION}"
    if [ "${DRY_RUN}" -eq 1 ]; then
        info "Would download and install godyl ${GODYL_VERSION}"
        info "Would install tools to ${INSTALL_DIR}"
        return 0
    fi

    tempfile="${TEMP_DIR}/install.sh"
    curl ${DISABLE_SSL:+-k} -sSL "https://raw.githubusercontent.com/idelchi/godyl/refs/heads/dev/install.sh" -o "${tempfile}"
    sh -s -- -v "${GODYL_VERSION}" -d "${TEMP_DIR}" ${DISABLE_SSL:+-k} < "${tempfile}"

    success "Installing tools to ${INSTALL_DIR}"

    # Install tools using godyl
    "${TEMP_DIR}/godyl" --output="${INSTALL_DIR}" - <<YAML
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

    success "All tools installed successfully to ${INSTALL_DIR}"
}

main() {
    parse_args "$@"

    # Check for required commands
    need_cmd curl
    need_cmd mktemp

    # Setup temporary directory
    setup_temp_dir

    # Create installation directory if it doesn't exist
    mkdir -p "${INSTALL_DIR}"

    # Install tools
    install_tools
}

main "$@"
