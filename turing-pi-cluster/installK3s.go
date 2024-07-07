package main

import (
	"fmt"

	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var k3sServerCommand *remote.Command
var k3sAgentCommand *remote.Command
var k3sServerCommandErr error
var k3sAgentCommandError error

func installK3s(ctx *pulumi.Context, nodes []node) error {

	for index, node := range nodes {
		if index == 0 {
			// We are on the first node, so we will install the K3s server
			k3sFirstNodeCommand := "curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION=\"v1.29.6+k3s1\" INSTALL_K3S_EXEC=\"--flannel-backend=none --disable-network-policy --write-kubeconfig-mode 644 --disable servicelb --token myrandompassword --disable-cloud-controller --disable local-storage --disable-kube-proxy  --disable traefik --data-dir /mnt/data/k3s\" sh -"

			conn := remote.ConnectionArgs{
				Host:       pulumi.String(node.addresss),
				User:       pulumi.String(node.username),
				PrivateKey: pulumi.String(node.sshKey),
			}

			k3sServerCommand, k3sServerCommandErr = remote.NewCommand(ctx, "InstallK3s"+node.addresss, &remote.CommandArgs{
				Connection: conn,
				Create:     pulumi.String(k3sFirstNodeCommand),
				Delete:     pulumi.String("sh /usr/local/bin/k3s-uninstall.sh"),
			}, pulumi.DependsOn([]pulumi.Resource{PartitioningScript[0], PartitioningScript[1], PartitioningScript[2]}))

			if k3sServerCommandErr != nil {
				return k3sServerCommandErr
			}

		} else {
			// We are on the other nodes, so we will install the K3s agents
			k3sSubsequentNodeCommand := fmt.Sprintf("curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION=\"v1.29.6+k3s1\" INSTALL_K3S_EXEC=\"agent --server https://%s:6443 --token myrandompassword --data-dir /mnt/data/k3s\" sh -s -", nodes[0].addresss)
			conn := remote.ConnectionArgs{
				Host:       pulumi.String(node.addresss),
				User:       pulumi.String(node.username),
				PrivateKey: pulumi.String(node.sshKey),
			}

			k3sAgentCommand, k3sAgentCommandError = remote.NewCommand(ctx, "InstallK3s"+node.addresss, &remote.CommandArgs{
				Connection: conn,
				Create:     pulumi.String(k3sSubsequentNodeCommand),
				Delete:     pulumi.String("sh /usr/local/bin/k3s-agent-uninstall.sh"),
			}, pulumi.DependsOn([]pulumi.Resource{PartitioningScript[0], PartitioningScript[1], PartitioningScript[2], k3sServerCommand}))

			if k3sAgentCommandError != nil {
				return k3sAgentCommandError
			}

		}
	}

	return nil
}
