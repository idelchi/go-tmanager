#!/bin/sh

# curl -sSL https://raw.githubusercontent.com/idelchi/godyl/refs/heads/dev/scripts/heredoc.sh | sh -s

dir=$(mktemp -d)

install_dir=${1:-~/.local/bin}
disable_ssl=${2:-0}

if [ "${disable_ssl}" -eq 1 ]; then
    flag="-k"
else
    flag=""
fi

curl ${flag} -sSL https://raw.githubusercontent.com/idelchi/godyl/refs/heads/dev/install.sh | sh -s -- -v v0.2-beta -d ${dir}

${dir}/godyl --output=${install_dir} - <<YAML
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

rm -rf ${dir}
