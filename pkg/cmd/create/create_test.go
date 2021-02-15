package create_test

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/jenkins-x-plugins/jx-registry/pkg/amazon/ecrs/fakeecr"
	"github.com/jenkins-x-plugins/jx-registry/pkg/cmd/create"
	jxcore "github.com/jenkins-x/jx-api/v4/pkg/apis/core/v4beta1"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	// generateTestOutput enable to regenerate the expected output
	generateTestOutput = false
)

func TestCreateForNonEKS(t *testing.T) {
	_, o := create.NewCmdCreate()

	o.Requirements = &jxcore.RequirementsConfig{
		Cluster: jxcore.ClusterConfig{
			Provider: "gke",
		},
	}

	err := o.Run()
	require.NoError(t, err, "failed to run")
}

func TestCreateForEKS(t *testing.T) {
	_, o := create.NewCmdCreate()

	o.Requirements = &jxcore.RequirementsConfig{
		Cluster: jxcore.ClusterConfig{
			Provider: "eks",
		},
	}

	o.AWSRegion = "dummy"
	o.Config = &aws.Config{}
	o.AppName = "myapp"
	o.ECRClient = fakeecr.NewFakeECR()

	err := o.Run()
	require.NoError(t, err, "failed to run")
}
