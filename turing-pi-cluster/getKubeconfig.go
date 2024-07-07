package main

import (
	"fmt"
	"os"

	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func getKubeconfig(ctx *pulumi.Context, serverNode node) error {

	conn := remote.ConnectionArgs{
		Host:       pulumi.String(serverNode.addresss),
		User:       pulumi.String(serverNode.username),
		PrivateKey: pulumi.String(serverNode.sshKey),
	}

	getKubeConfig, err := remote.NewCommand(ctx, "getKubeConfig", &remote.CommandArgs{
		Connection: conn,
		Create:     pulumi.String(fmt.Sprintf("cat /etc/rancher/k3s/k3s.yaml | sed 's/127.0.0.1/%s/g'", serverNode.addresss)),
		Delete:     pulumi.String("rm -f /$HOME/.kube/config"),
	}, pulumi.DependsOn([]pulumi.Resource{k3sServerCommand}))

	if err != nil {
		return err
	}

	getKubeConfig.Stdout.ApplyT(func(stdout string) error {
		return saveKubeconfig("/$HOME/.kube/config", stdout)

	})

	if err != nil {
		return err
	}

	return nil
}

// Helper function to write content to a file
func saveKubeconfig(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	return nil
}
