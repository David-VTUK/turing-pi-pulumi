package main

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apiextensions"
	helmv4 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v4"
	v1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func installCilium(ctx *pulumi.Context, ciliumVersion, k8sClusterPoolCidr, k8sServiceHost, k8sServicePort string) error {

	ciliumChart, err := helmv4.NewChart(ctx, "helm-chart-cilium", &helmv4.ChartArgs{
		Chart:   pulumi.String("cilium"),
		Version: pulumi.String(ciliumVersion),
		RepositoryOpts: &helmv4.RepositoryOptsArgs{
			Repo: pulumi.String("https://helm.cilium.io/"),
		},
		Namespace: pulumi.String("kube-system"),
		Values: pulumi.Map{
			"k8sServiceHost": pulumi.String(k8sServiceHost),
			"k8sServicePort": pulumi.String(k8sServicePort),

			"hubble": pulumi.Map{
				"enabled": pulumi.BoolPtr(true),
			},
			"ipam": pulumi.Map{
				"operator": pulumi.Map{
					"clusterPoolIPv4PodCIDRList": pulumi.StringArray{pulumi.String(k8sClusterPoolCidr)},
				},
			},
			"bgpControlPlane": pulumi.Map{
				"enabled": pulumi.BoolPtr(true),
			},
		},
	})

	if err != nil {
		return err
	}

	_, err = apiextensions.NewCustomResource(ctx, "bgp-peering-policy", &apiextensions.CustomResourceArgs{
		ApiVersion: pulumi.String("cilium.io/v2alpha1"),
		Kind:       pulumi.String("CiliumBGPClusterConfig"),
		Metadata: &v1.ObjectMetaArgs{
			Name: pulumi.String("default"),
		},
		OtherFields: kubernetes.UntypedArgs{
			"spec": pulumi.Map{
				"bgpInstances": pulumi.Array{
					pulumi.Map{
						"name":     pulumi.String("k8s-cluster"),
						"localASN": pulumi.Int(64512),
						"peers": pulumi.Array{
							pulumi.Map{
								"name":        pulumi.String("mikrotik"),
								"peerASN":     pulumi.Int(64512),
								"peerAddress": pulumi.String("172.16.10.1"),
								"peerConfigRef": pulumi.Map{
									"name": pulumi.String("mikrotik"),
								},
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

	_, err = apiextensions.NewCustomResource(ctx, "bgp-peer-config", &apiextensions.CustomResourceArgs{
		ApiVersion: pulumi.String("cilium.io/v2alpha1"),
		Kind:       pulumi.String("CiliumBGPPeerConfig"),
		Metadata: &v1.ObjectMetaArgs{
			Name: pulumi.String("mikrotik"),
		},
		OtherFields: kubernetes.UntypedArgs{
			"spec": pulumi.Map{
				"timers": pulumi.Map{
					"holdTimeSeconds":      pulumi.Int(9),
					"keepAliveTimeSeconds": pulumi.Int(3),
				},
				"gracefulRestart": pulumi.Map{
					"enabled":            pulumi.BoolPtr(true),
					"restartTimeSeconds": pulumi.Int(15),
				},
				"families": pulumi.Array{
					pulumi.Map{
						"afi":  pulumi.String("ipv4"),
						"safi": pulumi.String("unicast"),
						"advertisements": pulumi.Map{
							"matchLabels": pulumi.Map{
								"advertise": pulumi.String("bgp"),
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
				"blocks": pulumi.Array{
					pulumi.Map{
						"cidr": pulumi.String("10.200.200.0/24"),
					},
				},
			},
		},
	}, pulumi.DependsOn([]pulumi.Resource{ciliumChart}))

	if err != nil {
		return err
	}

	_, err = apiextensions.NewCustomResource(ctx, "bgp-advertisements", &apiextensions.CustomResourceArgs{
		ApiVersion: pulumi.String("cilium.io/v2alpha1"),
		Kind:       pulumi.String("CiliumBGPAdvertisement"),
		Metadata: &v1.ObjectMetaArgs{
			Name: pulumi.String("bgp-advertisements"),
			Labels: pulumi.StringMap{
				"advertise": pulumi.String("bgp"),
			},
		},
		OtherFields: kubernetes.UntypedArgs{
			"spec": pulumi.Map{
				"advertisements": pulumi.Array{
					pulumi.Map{
						"advertisementType": pulumi.String("Service"),
						"service": pulumi.Map{
							"addresses": pulumi.StringArray{
								pulumi.String("LoadBalancerIP"),
							},
						},
						"selector": pulumi.Map{
							"matchExpressions": pulumi.Array{
								pulumi.Map{
									"key":      pulumi.String("somekey"),
									"operator": pulumi.String("NotIn"),
									"values":   pulumi.StringArray{pulumi.String("never-used-value")},
								},
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

	return nil
}
