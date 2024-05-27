package main

import (
	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func prepNode(ctx *pulumi.Context, address, username, password string) error {

	_ = remote.ConnectionArgs{
		AgentSocketPath:    nil,
		DialErrorLimit:     nil,
		Host:               nil,
		Password:           nil,
		PerDialTimeout:     nil,
		Port:               nil,
		PrivateKey:         nil,
		PrivateKeyPassword: nil,
		Proxy:              nil,
		User:               nil,
	}
	return nil
}
