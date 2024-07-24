package main

import (
	helmv4 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v4"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

// cilium install --wait --version 1.15.6 --set=ipam.operator.clusterPoolIPv4PodCIDRList="10.42.0.0/16" --set=k8sServiceHost="172.16.10.224" --set=k8sServicePort="6443"

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Get values from config
		conf := config.New(ctx, "")

		// Get Node IP's
		ciliumVersion := conf.Require("cilium-version")

		_, err := helmv4.NewChart(ctx, "helm-chart-cilium", &helmv4.ChartArgs{
			Chart:   pulumi.String("cilium"),
			Version: pulumi.String(ciliumVersion),
			RepositoryOpts: &helmv4.RepositoryOptsArgs{
				Repo: pulumi.String("https://helm.cilium.io/"),
			},
			Namespace: pulumi.String("kube-system"),
			Values: pulumi.Map{
				"hubble": pulumi.Map{
					"enabled": pulumi.String("true"),
				},
				"ipam": pulumi.Map{
					"operator": pulumi.Map{
						"clusterPoolIPv4PodCIDRList": pulumi.StringArray{pulumi.String("10.42.0.0/16")},
					},
				},
			},
		})

		if err != nil {
			return err
		}

		return nil
	})

}
