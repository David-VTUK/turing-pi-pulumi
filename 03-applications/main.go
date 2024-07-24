package main

import (
	helmv4 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v4"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		_, err := helmv4.NewChart(ctx, "helm-chart-longhorn", &helmv4.ChartArgs{
			Chart: pulumi.String("longhorn"),
			RepositoryOpts: &helmv4.RepositoryOptsArgs{
				Repo: pulumi.String("https://charts.longhorn.io"),
			},
		})

		if err != nil {
			return err
		}

		return nil
	})
}
