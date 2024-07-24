package main

import (
	"fmt"

	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var k3sServerCommand *remote.Command
var k3sAgentCommand *remote.Command
var err error

func installK3s(ctx *pulumi.Context, nodes []Node, k3sversion string) error {

	for index, node := range nodes {
		if index == 0 {
			// We are on the first node, so we will install the K3s server
			k3sFirstNodeCommand := fmt.Sprintf("curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION=\"%s\" INSTALL_K3S_EXEC=\"--flannel-backend=none --disable-network-policy --write-kubeconfig-mode 644 --disable servicelb --token myrandompassword --disable-cloud-controller --disable local-storage --disable-kube-proxy  --disable traefik\" sh -", k3sversion)

			conn := remote.ConnectionArgs{
				Host:       pulumi.String(node.addresss),
				User:       pulumi.String(node.username),
				PrivateKey: pulumi.String(node.sshKey),
			}

			k3sServerCommand, err = remote.NewCommand(ctx, "InstallK3s"+node.addresss, &remote.CommandArgs{
				Connection: conn,
				Create:     pulumi.String(k3sFirstNodeCommand),
				Delete:     pulumi.String("sh /usr/local/bin/k3s-uninstall.sh"),
			})

			if err != nil {
				return err
			}

		} else {
			// We are on the other nodes, so we will install the K3s agents
			k3sSubsequentNodeCommand := fmt.Sprintf("curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION=\"%s\" INSTALL_K3S_EXEC=\"agent --server https://%s:6443 --token myrandompassword\" sh -s -", k3sversion, nodes[0].addresss)
			conn := remote.ConnectionArgs{
				Host:       pulumi.String(node.addresss),
				User:       pulumi.String(node.username),
				PrivateKey: pulumi.String(node.sshKey),
			}

			k3sAgentCommand, err = remote.NewCommand(ctx, "InstallK3s"+node.addresss, &remote.CommandArgs{
				Connection: conn,
				Create:     pulumi.String(k3sSubsequentNodeCommand),
				Delete:     pulumi.String("sh /usr/local/bin/k3s-agent-uninstall.sh"),
			}, pulumi.DependsOn([]pulumi.Resource{k3sServerCommand}))

			if err != nil {
				return err
			}

		}
	}

	return nil
}
