#!/bin/sh

# curl -sSL https://raw.githubusercontent.com/idelchi/godyl/refs/heads/dev/heredoc.sh | sh -s

dir=$(mktemp -d)

install_dir=${1:-~/.local/bin}

curl -sSL https://raw.githubusercontent.com/idelchi/godyl/refs/heads/main/scripts/install.sh | sh -s -- -v v0.1-beta -o "${dir}"

"${dir}"/godyl --output="${install_dir}" - <<YAML
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

rm -rf "${dir}"
