package create

import (
	"context"
	"fmt"
	"github.com/jenkins-x-plugins/jx-registry/pkg/amazon/ecr"
	"github.com/jenkins-x-plugins/jx-registry/pkg/rootcmd"
	jxcore "github.com/jenkins-x/jx-api/v4/pkg/apis/core/v4beta1"
	"github.com/jenkins-x/jx-api/v4/pkg/client/clientset/versioned"
	"github.com/jenkins-x/jx-gitops/pkg/variablefinders"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cmdrunner"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras/helper"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras/templates"
	"github.com/jenkins-x/jx-helpers/v3/pkg/gitclient"
	"github.com/jenkins-x/jx-helpers/v3/pkg/gitclient/cli"
	"github.com/jenkins-x/jx-helpers/v3/pkg/kube/jxclient"
	"github.com/jenkins-x/jx-helpers/v3/pkg/options"
	"github.com/jenkins-x/jx-helpers/v3/pkg/termcolor"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"github.com/pkg/errors"
	"github.com/sethvargo/go-envconfig"
	"github.com/spf13/cobra"
	"strings"
)

var (
	info = termcolor.ColorInfo

	cmdLong = templates.LongDesc(`
		Lazy create a container registry for ECR
`)

	cmdExample = templates.Examples(`
		# lets ensure we have an ECR registry setup
		%s create
	`)
)

// Options the options for this command
type Options struct {
	options.BaseOptions
	AWSProfile           string `env:"AWS_PROFILE"`
	AWSRegion            string `env:"AWS_REGION"`
	Registry             string `env:"DOCKER_REGISTRY"`
	RegistryOrganisation string `env:"DOCKER_REGISTRY_ORG"`
	AppName              string `env:"APP_NAME"`
	ECRSuffix            string
	Namespace            string
	JXClient             versioned.Interface
	GitClient            gitclient.Interface
	CommandRunner        cmdrunner.CommandRunner
	Requirements         *jxcore.RequirementsConfig
}

// NewCmdCreate creates a command object for the command
func NewCmdCreate() (*cobra.Command, *Options) {
	o := &Options{}

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Lazy create a container registry for ECR",
		Long:    cmdLong,
		Example: fmt.Sprintf(cmdExample, rootcmd.BinaryName),
		Run: func(cmd *cobra.Command, args []string) {
			err := o.Run()
			helper.CheckErr(err)
		},
	}

	if o.Ctx == nil {
		o.Ctx = cmd.Context()
	}
	if o.Ctx == nil {
		o.Ctx = context.TODO()
	}
	err := envconfig.Process(o.Ctx, o)
	if err != nil {
		log.Logger().Warnf("failed to default env vars: %s", err.Error())
	}

	cmd.Flags().StringVarP(&o.Namespace, "namespace", "n", "", "The namespace. Defaults to the current namespace")
	cmd.Flags().StringVarP(&o.AWSProfile, "aws-profile", "", o.AWSProfile, "The AWS profile to use. Defaults to $AWS_PROFILE")
	cmd.Flags().StringVarP(&o.AWSRegion, "aws-region", "", o.AWSRegion, "The AWS region. Defaults to $AWS_REGION or its read from the 'jx-requirements.yml' for the development environment")
	cmd.Flags().StringVarP(&o.Registry, "registry", "r", o.Registry, "The registry to use. Defaults to $DOCKER_REGISTRY")
	cmd.Flags().StringVarP(&o.RegistryOrganisation, "organisation", "o", o.RegistryOrganisation, "The registry organisation to use. Defaults to $DOCKER_REGISTRY_ORG")
	cmd.Flags().StringVarP(&o.AppName, "app", "a", o.AppName, "The app name to use. Defaults to $APP_NAME")
	cmd.Flags().StringVarP(&o.ECRSuffix, "ecr-registry-suffix", "", ".amazonaws.com", "The registry suffix to check if we are using ECR")

	o.BaseOptions.AddBaseFlags(cmd)
	return cmd, o
}

func (o *Options) Validate() error {
	if o.GitClient == nil {
		o.GitClient = cli.NewCLIClient("", o.CommandRunner)
	}
	var err error
	o.JXClient, o.Namespace, err = jxclient.LazyCreateJXClientAndNamespace(o.JXClient, o.Namespace)
	if err != nil {
		return errors.Wrapf(err, "failed to create jxClient")
	}

	o.Requirements, err = variablefinders.FindRequirements(o.GitClient, o.JXClient, o.Namespace, "")
	if err != nil {
		return errors.Wrapf(err, "failed to load requirements from dev environment")
	}
	if o.Requirements == nil {
		return errors.Errorf("no requirements found for dev environment")
	}

	if o.Requirements.Cluster.Provider == "eks" {
		if o.AWSRegion == "" {
			o.AWSRegion = o.Requirements.Cluster.Region

			if o.AWSRegion == "" {
				log.Logger().Warnf("could not find the AWS region in the 'jx-requirements.yml' file in cluster.region or in $AWS_REGION")
				return options.MissingOption("aws-region")
			}
		}
	}
	return nil

}
func (o *Options) Run() error {
	err := o.Validate()
	if err != nil {
		return errors.Wrapf(err, "failed to validate options")
	}
	if o.Requirements.Cluster.Provider != "eks" {
		return nil
	}
	registry := o.Requirements.Cluster.Registry
	if registry != "" && strings.HasSuffix(registry, o.ECRSuffix) {
		return nil
	}

	log.Logger().Infof("verifying that container registry %s with organisation %s and app name %s has an ECR associated with it", info(registry), info(o.RegistryOrganisation), info(o.AppName))
	err = ecr.LazyCreateRegistry(o.AWSProfile, o.AWSRegion, o.Registry, o.RegistryOrganisation, o.AppName)
	if err != nil {
		return errors.Wrapf(err, "failed to lazy create the ECR registry")
	}
	return nil
}
