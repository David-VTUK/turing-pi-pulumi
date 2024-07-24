package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Get values from config
		conf := config.New(ctx, "")

		// Get Cilium Version
		ciliumVersion := conf.Require("cilium-version")

		err := installCilium(ctx, ciliumVersion)
		if err != nil {
			return err

		}

		return nil
	})

}
