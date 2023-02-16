package main

import (
	"fmt"
	// Pulumi Google Classic Providers
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/projects"
	// Pulumi Google Native Providers
	run "github.com/pulumi/pulumi-google-native/sdk/go/google/run/v2"
	// Pulumi Core
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

// Enabled Google Cloud Services
var services = []string{
	"compute.googleapis.com",
	"run.googleapis.com",
}

// Dependencies
var dependencies []pulumi.Resource

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Load Pulumi Configuration
		conf := config.New(ctx, "google-native")
		project := conf.Require("project")

		// Enable all Required Google Cloud Project Service API's

		for idx, service := range services {
			newService, err := projects.NewService(ctx, fmt.Sprintf("project-service-%d", idx), &projects.ServiceArgs{
				DisableDependentServices: pulumi.Bool(true),
				Service:                  pulumi.String(service),
				Project:                  pulumi.String(project),
			})
			if err != nil {
				return err
			}
			dependencies = append(dependencies, newService)
		}

		// Create Cloud Run Service with public Hello World container image running on Port 8080 in europe-west1
		run.NewService(ctx, "cloud-run-service", &run.ServiceArgs{
			Project:     pulumi.String(project),
			ServiceId:   pulumi.String("cloud-run-service-001"),
			Description: pulumi.String("Cloud Run Service, Container: 'Hello World', Region: europe-west1, Auth: Public/NoAuth"),
			Location:    pulumi.String("europe-west1"),
			Template: &run.GoogleCloudRunV2RevisionTemplateArgs{
				Containers: &run.GoogleCloudRunV2ContainerArray{
					&run.GoogleCloudRunV2ContainerArgs{
						Image: pulumi.String("gcr.io/cloudrun/hello"),
						Ports: &run.GoogleCloudRunV2ContainerPortArray{
							&run.GoogleCloudRunV2ContainerPortArgs{
								Name:          pulumi.String("http1"),
								ContainerPort: pulumi.Int(8080),
							},
						},
					},
				},
			},
		}, pulumi.DependsOn(dependencies))


		run.NewServiceIamPolicy(ctx,  "cloud-run-service-iam-policy-no-auth", &run.ServiceIamPolicyArgs{
			
		})

		return nil
	})
}
