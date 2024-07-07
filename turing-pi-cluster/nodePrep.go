package main

import (
	"log"
	"os"

	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var PartitioningScript []*remote.Command

func prepNode(ctx *pulumi.Context, nodes []node) error {

	for _, v := range nodes {

		createPartitionScript, err := os.ReadFile("./assets/createPartition.sh")
		if err != nil {
			log.Fatalf("failed reading file: %s", err)
		}

		deletePartitionScript, err := os.ReadFile("./assets/deletePartition.sh")
		if err != nil {
			log.Fatalf("failed reading file: %s", err)
		}

		conn := remote.ConnectionArgs{
			Host:       pulumi.String(v.addresss),
			User:       pulumi.String(v.username),
			PrivateKey: pulumi.String(v.sshKey),
		}

		partitioningScript, err := remote.NewCommand(ctx, "configurePartitions"+v.addresss, &remote.CommandArgs{
			Connection: conn,
			Create:     pulumi.String(createPartitionScript),
			Delete:     pulumi.String(deletePartitionScript),
		})

		PartitioningScript = append(PartitioningScript, partitioningScript)

		if err != nil {
			return err
		}
	}

	return nil

}
