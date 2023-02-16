package main

import (
	"fmt"

	"github.com/google/uuid"

	// Pulumi Google Native Providers
	crm "github.com/pulumi/pulumi-google-native/sdk/go/google/cloudresourcemanager/v3"
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
		parentId := conf.Require("parentId")
		projectId, err := conf.TryInt("projectId")
		if err != nil {
			projectId = uuid.New()
		}
		projectName, err := conf.TryInt("projectName")
		if err != nil {
			projectName = projectId
		}

		project := crm.Project.Project
		project, err = crm.Project.NewProject(ctx, fmt.Sprintf("google-cloud-project-%s", projectId), &crm.Project.ProjectArgs{
			Parent:      pulumi.String(parentId),
			ProjectId:   pulumi.String(projectId),
			DisplayName: pulumi.String(projectName),
		})
		if err != nil {
			return err
		}

		// Export Project
		ctx.Export(fmt.Sprintf("google-cloud-project-%s", projectId), project)
		// Export ProjectId
		ctx.Export(fmt.Sprintf("google-cloud-project-id-%s", projectId), projectId)

		return nil
	})
}
