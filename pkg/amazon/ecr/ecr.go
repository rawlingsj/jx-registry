package ecr

import (
	"fmt"
	"github.com/jenkins-x-plugins/jx-registry/pkg/amazon/session"
	"github.com/jenkins-x/jx-helpers/v3/pkg/options"
	"github.com/jenkins-x/jx-helpers/v3/pkg/termcolor"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/sts"
)

// GetAccountIDAndRegion returns the current account ID and region
func GetAccountIDAndRegion(profile string, region string) (string, string, error) {
	sess, err := session.NewAwsSession(profile, region)
	// We nee to get the region from the connected cluster instead of the one configured for the calling user
	// as it might not be found and it would then use the default (us-west-2)
	_, region, err = session.GetCurrentlyConnectedRegionAndClusterName()
	if err != nil {
		return "", "", err
	}
	svc := sts.New(sess)

	input := &sts.GetCallerIdentityInput{}

	result, err := svc.GetCallerIdentity(input)
	if err != nil {
		return "", region, err
	}
	if result.Account != nil {
		return *result.Account, region, nil
	}
	return "", region, fmt.Errorf("Could not find the AWS Account ID!")
}

// GetContainerRegistryHost
func GetContainerRegistryHost() (string, error) {
	accountId, region, err := GetAccountIDAndRegion("", "")
	if err != nil {
		return "", err
	}
	return accountId + ".dkr.ecr." + region + ".amazonaws.com", nil
}

// LazyCreateRegistry lazily creates the ECR registry if it does not already exist
func LazyCreateRegistry(profileOption string, region string, dockerRegistry string, orgName string, appName string) error {
	if region == "" {
		return options.MissingOption("aws-region")
	}

	// strip any tag/version from the app name
	idx := strings.Index(appName, ":")
	if idx > 0 {
		appName = appName[0:idx]
	}
	repoName := appName
	if orgName != "" {
		repoName = orgName + "/" + appName
	}
	repoName = strings.ToLower(repoName)
	log.Logger().Infof("Let's ensure that we have an ECR repository for the image %s", termcolor.ColorInfo(repoName))
	sess, err := session.NewAwsSession(profileOption, region)
	if err != nil {
		return err
	}
	svc := ecr.New(sess)
	repoInput := &ecr.DescribeRepositoriesInput{
		RepositoryNames: []*string{
			aws.String(repoName),
		},
	}
	result, err := svc.DescribeRepositories(repoInput)
	if aerr, ok := err.(awserr.Error); !ok || aerr.Code() != ecr.ErrCodeRepositoryNotFoundException {
		return err
	}
	for _, repo := range result.Repositories {
		name := repo.String()
		log.Logger().Infof("Found repository: %s", name)
		if name == repoName {
			return nil
		}
	}
	createRepoInput := &ecr.CreateRepositoryInput{
		RepositoryName: aws.String(repoName),
	}
	createResult, err := svc.CreateRepository(createRepoInput)
	if err != nil {
		return fmt.Errorf("Failed to create the ECR repository for %s due to: %s", repoName, err)
	}
	repo := createResult.Repository
	if repo != nil {
		u := repo.RepositoryUri
		if u != nil {
			if !strings.HasPrefix(*u, dockerRegistry) {
				log.Logger().Warnf("Created ECR repository (%s) doesn't match registry configured for team (%s)",
					termcolor.ColorInfo(*u), termcolor.ColorInfo(dockerRegistry))
			} else {
				log.Logger().Infof("Created ECR repository: %s", termcolor.ColorInfo(*u))
			}
		}
	}
	return nil
}
