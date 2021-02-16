package ecrs

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/jenkins-x-plugins/jx-registry/pkg/amazon"
	"github.com/jenkins-x/jx-helpers/v3/pkg/options"
	"github.com/jenkins-x/jx-helpers/v3/pkg/termcolor"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"github.com/pkg/errors"
	"github.com/sethvargo/go-envconfig"
	"github.com/spf13/cobra"
)

type ECRClient interface {
	DescribeRepositories(context.Context, *ecr.DescribeRepositoriesInput, ...func(*ecr.Options)) (*ecr.DescribeRepositoriesOutput, error)
	CreateRepository(ctx context.Context, params *ecr.CreateRepositoryInput, optFns ...func(*ecr.Options)) (*ecr.CreateRepositoryOutput, error)
}

type Options struct {
	amazon.Options
	Registry             string `env:"DOCKER_REGISTRY"`
	RegistryOrganisation string `env:"DOCKER_REGISTRY_ORG"`
	AppName              string `env:"APP_NAME"`
	ECRClient            ECRClient
}

func (o *Options) AddFlags(cmd *cobra.Command) {
	o.Options.AddFlags(cmd)

	cmd.Flags().StringVarP(&o.Registry, "registry", "r", o.Registry, "The registry to use. Defaults to $DOCKER_REGISTRY")
	cmd.Flags().StringVarP(&o.RegistryOrganisation, "organisation", "o", o.RegistryOrganisation, "The registry organisation to use. Defaults to $DOCKER_REGISTRY_ORG")
	cmd.Flags().StringVarP(&o.AppName, "app", "a", o.AppName, "The app name to use. Defaults to $APP_NAME")
}

func (o *Options) Validate() error {
	cfg, err := o.GetConfig()
	if err != nil {
		return errors.Wrapf(err, "failed to create AWS config")
	}
	if cfg == nil {
		return errors.Errorf("no AWS config")
	}
	return nil
}

// EnvProcess processes the environment variable defaults
func (o *Options) EnvProcess() {
	err := envconfig.Process(o.GetContext(), o)
	if err != nil {
		log.Logger().Warnf("failed to default env vars: %s", err.Error())
	}
}

// LazyCreateRegistry lazily creates the ECR registry if it does not already exist
func (o *Options) LazyCreateRegistry(appName string) error {
	ctx := o.GetContext()
	cfg, err := o.GetConfig()
	if err != nil {
		return errors.Wrapf(err, "failed to create the AWS configuration")
	}

	region := o.AWSRegion
	if region == "" {
		return options.MissingOption("aws-region")
	}

	// strip any tag/version from the app name
	idx := strings.Index(appName, ":")
	if idx > 0 {
		appName = appName[0:idx]
	}
	repoName := appName
	if o.RegistryOrganisation != "" {
		repoName = o.RegistryOrganisation + "/" + appName
	}
	repoName = strings.ToLower(repoName)
	log.Logger().Infof("Let's ensure that we have an ECR repository for the image %s", termcolor.ColorInfo(repoName))

	if o.ECRClient == nil {
		o.ECRClient = ecr.NewFromConfig(*cfg, nil)
	}
	svc := o.ECRClient

	repoInput := &ecr.DescribeRepositoriesInput{
		RepositoryNames: []string{repoName},
	}
	result, err := svc.DescribeRepositories(ctx, repoInput)
	if err != nil {
		if _, ok := err.(*types.RepositoryNotFoundException); !ok {
			return errors.Wrapf(err, "failed to check for repository")
		}
	}
	if result != nil {
		for _, repo := range result.Repositories {
			if repo.RepositoryName == nil {
				continue
			}
			name := *repo.RepositoryName
			log.Logger().Infof("Found repository: %s", name)
			if name == repoName {
				return nil
			}
		}
	}
	createRepoInput := &ecr.CreateRepositoryInput{
		RepositoryName: aws.String(repoName),
	}
	createResult, err := svc.CreateRepository(ctx, createRepoInput)
	if err != nil {
		return fmt.Errorf("Failed to create the ECR repository for %s due to: %s", repoName, err)
	}
	repo := createResult.Repository
	if repo != nil {
		u := repo.RepositoryUri
		if u != nil {
			log.Logger().Infof("Created ECR repository: %s", termcolor.ColorInfo(*u))
		}
	}
	return nil
}
