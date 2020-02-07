// Copyright 2020 The Terraformer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aws

import (
	"context"

	"github.com/GoogleCloudPlatform/terraformer/terraform_utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codebuild"
)

var codebuildAllowEmptyValues = []string{"tags."}

type CodeBuildGenerator struct {
	AWSService
}

func (g *CodeBuildGenerator) loadProjects(svc *codebuild.Client) error {
	output, err := svc.ListProjectsRequest(&codebuild.ListProjectsInput{}).Send(context.Background())
	if err != nil {
		return err
	}
	for _, project := range output.Projects {
		g.Resources = append(g.Resources, terraform_utils.NewSimpleResource(
			project,
			project,
			"aws_codebuild_project",
			"aws",
			codebuildAllowEmptyValues))
	}
	return nil
}

func (g *CodeBuildGenerator) loadSourceCredentials(svc *codebuild.Client) error {
	output, err := svc.ListSourceCredentialsRequest(&codebuild.ListSourceCredentialsInput{}).Send(context.Background())
	if err != nil {
		return err
	}
	for _, sourceCredentials := range output.SourceCredentialsInfos {
		resourceArn := aws.StringValue(sourceCredentials.Arn)
		g.Resources = append(g.Resources, terraform_utils.NewSimpleResource(
			resourceArn,
			resourceArn,
			"aws_codebuild_source_credential",
			"aws",
			codebuildAllowEmptyValues))
	}
	return nil
}

func (g *CodeBuildGenerator) InitResources() error {
	config, e := g.generateConfig()
	if e != nil {
		return e
	}
	svc := codebuild.New(config)

	if err := g.loadProjects(svc); err != nil {
		return err
	}
	if err := g.loadSourceCredentials(svc); err != nil {
		return err
	}

	return nil
}
