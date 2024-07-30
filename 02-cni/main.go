package main

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apiextensions"
	helmv4 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v4"
	v1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func installCilium(ctx *pulumi.Context, ciliumVersion string) error {

	ciliumChart, err := helmv4.NewChart(ctx, "helm-chart-cilium", &helmv4.ChartArgs{
		Chart:   pulumi.String("cilium"),
		Version: pulumi.String(ciliumVersion),
		RepositoryOpts: &helmv4.RepositoryOptsArgs{
			Repo: pulumi.String("https://helm.cilium.io/"),
		},
		Namespace: pulumi.String("kube-system"),
		Values: pulumi.Map{
			"k8sServiceHost": pulumi.String("172.16.10.220"),
			"k8sServicePort": pulumi.String("6443"),
			//			"annotations": pulumi.Map{
			//				"meta.helm.sh/release-name":      pulumi.String("helm-chart-cilium"),
			//				"meta.helm.sh/release-namespace": pulumi.String("kube-system"),
			//			},

			"hubble": pulumi.Map{
				"enabled": pulumi.String("true"),
			},
			"ipam": pulumi.Map{
				"operator": pulumi.Map{
					"clusterPoolIPv4PodCIDRList": pulumi.StringArray{pulumi.String("10.42.0.0/16")},
				},
			},
			"bgpControlPlane": pulumi.Map{
				"enabled": pulumi.String("true"),
			},
		},
	})

	if err != nil {
		return err
	}

	// Set up BGP Peering with Mikrotik Router
	_, err = apiextensions.NewCustomResource(ctx, "bgp-peering-policy", &apiextensions.CustomResourceArgs{
		ApiVersion: pulumi.String("cilium.io/v2alpha1"),
		Kind:       pulumi.String("CiliumBGPPeeringPolicy"),
		Metadata: &v1.ObjectMetaArgs{
			Name: pulumi.String("default"),
		},
		OtherFields: kubernetes.UntypedArgs{
			"spec": pulumi.Map{
				"virtualRouters": pulumi.Array{
					pulumi.Map{
						"localASN": pulumi.Int(64511),
						"neighbors": pulumi.Array{
							pulumi.Map{
								"peerAddress": pulumi.String("172.16.10.1/32"),
								"peerASN":     pulumi.Int(64512),
							},
						},
					},
				},
			},
		},
	}, pulumi.DependsOn([]pulumi.Resource{ciliumChart}))

	if err != nil {
		return err
	}

	_, err = apiextensions.NewCustomResource(ctx, "bgp-ip-pool", &apiextensions.CustomResourceArgs{
		ApiVersion: pulumi.String("cilium.io/v2alpha1"),
		Kind:       pulumi.String("CiliumLoadBalancerIPPool"),
		Metadata: &v1.ObjectMetaArgs{
			Name: pulumi.String("default"),
		},
		OtherFields: kubernetes.UntypedArgs{
			"spec": pulumi.Map{
				"cidrs": pulumi.Array{
					pulumi.Map{
						"cidr": pulumi.String("10.50.20.0/24"),
					},
				},
			},
		},
	}, pulumi.DependsOn([]pulumi.Resource{ciliumChart}))

	// Set up Routes to Advertise

	if err != nil {
		return err
	}

	return nil
}
