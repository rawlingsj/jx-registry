module github.com/jenkins-x-plugins/jx-registry

require (
	github.com/aws/aws-sdk-go-v2 v1.2.0
	github.com/aws/aws-sdk-go-v2/config v1.1.1
	github.com/aws/aws-sdk-go-v2/service/ecr v1.1.1
	github.com/aws/smithy-go v1.1.0
	github.com/jenkins-x/jx-api/v4 v4.0.25
	github.com/jenkins-x/jx-gitops v0.2.35
	github.com/jenkins-x/jx-helpers/v3 v3.0.92
	github.com/jenkins-x/jx-logging/v3 v3.0.3
	github.com/pkg/errors v0.9.1
	github.com/sethvargo/go-envconfig v0.3.2
	github.com/spf13/cobra v1.1.1
	github.com/stretchr/testify v1.6.1
)

replace (
	k8s.io/api => k8s.io/api v0.20.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.20.2
	k8s.io/client-go => k8s.io/client-go v0.20.2
)

go 1.15
