FROM ubuntu:latest

# Install required dependencies
RUN apt-get update && \
    apt-get install -y curl git wget unzip jq nano


ARG GOVERSION=go1.22.3

# Install a specific version of Go
RUN curl -LO https://golang.org/dl/${GOVERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf ${GOVERSION}.linux-amd64.tar.gz && \
    rm ${GOVERSION}.linux-amd64.tar.gz

# Add Go binary to PATH
ENV PATH="/usr/local/go/bin:${PATH}"    

# Download and install Pulumi binary
RUN curl -fsSL https://get.pulumi.com | sh

# Add Pulumi binary to PATH
ENV PATH="/root/.pulumi/bin:${PATH}"

# Download and install kubectl binary
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
    chmod +x kubectl && \
    mv kubectl /usr/local/bin/


# Download and install the Cillium CLI
RUN CILIUM_CLI_VERSION=$(curl -s https://raw.githubusercontent.com/cilium/cilium-cli/main/stable.txt) && \
    CLI_ARCH=amd64 && \
    if [ "$(uname -m)" = "aarch64" ]; then CLI_ARCH=arm64; fi && \
    curl -L --fail --remote-name-all https://github.com/cilium/cilium-cli/releases/download/${CILIUM_CLI_VERSION}/cilium-linux-${CLI_ARCH}.tar.gz{,.sha256sum} && \
    sha256sum --check cilium-linux-${CLI_ARCH}.tar.gz.sha256sum && \
    tar xzvfC cilium-linux-${CLI_ARCH}.tar.gz /usr/local/bin 