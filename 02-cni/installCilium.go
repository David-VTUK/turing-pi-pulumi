package main

import (
	helmv4 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v4"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func installCilium(ctx *pulumi.Context, ciliumVersion string) error {

	_, err := helmv4.NewChart(ctx, "helm-chart-cilium", &helmv4.ChartArgs{
		Chart:   pulumi.String("cilium"),
		Version: pulumi.String(ciliumVersion),
		RepositoryOpts: &helmv4.RepositoryOptsArgs{
			Repo: pulumi.String("https://helm.cilium.io/"),
		},
		Namespace: pulumi.String("kube-system"),
		Values: pulumi.Map{
			"k8sServiceHost": pulumi.String("172.16.10.220"),
			"k8sServicePort": pulumi.String("6443"),
			"annotations": pulumi.Map{
				"meta.helm.sh/release-name": pulumi.String("helm-chart-cilium"),
				"meta.helm.sh/release-namespace": pulumi.String("kube-system"),
			},
			
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
}
