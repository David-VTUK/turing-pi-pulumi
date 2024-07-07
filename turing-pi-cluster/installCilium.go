package main

import (
	"log"
	"os"

	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func installCilium(ctx *pulumi.Context, serverNode node) error {

	installCiliumScript, err := os.ReadFile("./assets/installCilium.sh")
	if err != nil {
		log.Fatalf("failed reading file: %s", err)
	}

	conn := remote.ConnectionArgs{
		Host:       pulumi.String(serverNode.addresss),
		User:       pulumi.String(serverNode.username),
		PrivateKey: pulumi.String(serverNode.sshKey),
	}

	_, err = remote.NewCommand(ctx, "InstallCilium", &remote.CommandArgs{
		Connection: conn,
		Create:     pulumi.String(installCiliumScript),
	}, pulumi.DependsOn([]pulumi.Resource{k3sServerCommand,k3sAgentCommand}))

	if err != nil {
		return err
	}

	return nil
}

