package create_test

import (
	"github.com/jenkins-x-plugins/jx-registry/pkg/cmd/create"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	// generateTestOutput enable to regenerate the expected output
	generateTestOutput = false
)

func TestCreate(t *testing.T) {
	_, o := create.NewCmdCreate()

	require.NotNil(t, o, "options")
}
