#!/bin/bash

# Define variables
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml

# Install Cilium CLI
CILIUM_CLI_VERSION=$(curl -s https://raw.githubusercontent.com/cilium/cilium-cli/main/stable.txt)
CLI_ARCH=amd64
if [ "$(uname -m)" = "aarch64" ]; then CLI_ARCH=arm64; fi
curl -L --fail --remote-name-all https://github.com/cilium/cilium-cli/releases/download/${CILIUM_CLI_VERSION}/cilium-linux-${CLI_ARCH}.tar.gz{,.sha256sum}
sha256sum --check cilium-linux-${CLI_ARCH}.tar.gz.sha256sum
sudo tar xzvfC cilium-linux-${CLI_ARCH}.tar.gz /usr/local/bin
rm cilium-linux-${CLI_ARCH}.tar.gz{,.sha256sum}

# Install Cilium
cilium install --wait --version 1.15.6 --set=ipam.operator.clusterPoolIPv4PodCIDRList="10.42.0.0/16" --set=k8sServiceHost="172.16.10.224" --set=k8sServicePort="6443"


cilium hubble enable
