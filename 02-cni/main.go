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
		k8sClusterPoolCidr := conf.Require("k8s-cluster-pool-cidr")
  		k8sServiceHost := conf.Require("k8s-service-host")
  		k8sServicePort := conf.Require("k8s-service-port")


		err := installCilium(ctx, ciliumVersion, k8sClusterPoolCidr, k8sServiceHost, k8sServicePort)
		if err != nil {
			return err

		}

		return nil
	})

}
