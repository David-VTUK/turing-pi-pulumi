package main

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

type node struct {
	addresss string
	sshKey   string
	username string
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Get values from config
		conf := config.New(ctx, "")

		// Get Node IP's
		node1Address := conf.Require("node-1-ip")
		node2Address := conf.Require("node-2-ip")
		node3Address := conf.Require("node-3-ip")
		node4Address := conf.Require("node-4-ip")

		// Get K3s Version
		k3sVersion := conf.Require("k3s-version")

		// Get Node SSH Values
		username := conf.Require("ssh-user")
		sshKey := conf.RequireSecret("ssh-key")

		sshKey.ApplyT(func(sshKey string) error {

			nodes := []node{
				{addresss: node1Address, sshKey: sshKey, username: username},
				{addresss: node2Address, sshKey: sshKey, username: username},
				{addresss: node3Address, sshKey: sshKey, username: username},
				{addresss: node4Address, sshKey: sshKey, username: username},
			}

			err := installK3s(ctx, nodes, k3sVersion)
			if err != nil {
				return fmt.Errorf("failed to install K3S: %w", err)
			}

			err = getKubeconfig(ctx, nodes[0])
			if err != nil {
				return fmt.Errorf("failed to get Kubeconfig: %w", err)
			}

			err = installCilium(ctx, nodes[0])
			if err != nil {
				return fmt.Errorf("failed to install Cilium: %w", err)
			}

			return nil
		})

		return nil
	})
}
