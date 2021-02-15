package fakeecr

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/aws/smithy-go/middleware"
	"github.com/pkg/errors"
)

// FakeECR a fake ECR implementation for testing
type FakeECR struct {
	Repositories map[string]*types.Repository
}

func (f *FakeECR) DescribeRepositories(ctx context.Context, input *ecr.DescribeRepositoriesInput, opts ...func(*ecr.Options)) (*ecr.DescribeRepositoriesOutput, error) {
	var repos []types.Repository
	{
	}
	if input != nil && f.Repositories != nil {
		for _, name := range input.RepositoryNames {
			r := f.Repositories[name]
			if r != nil {
				repos = append(repos, *r)
			}
		}
	}
	return &ecr.DescribeRepositoriesOutput{
		Repositories:   repos,
		ResultMetadata: middleware.Metadata{},
	}, nil
}

func (f *FakeECR) CreateRepository(ctx context.Context, params *ecr.CreateRepositoryInput, opts ...func(*ecr.Options)) (*ecr.CreateRepositoryOutput, error) {
	if params.RepositoryName == nil {
		return nil, errors.Errorf("missing params.RepositoryName")
	}
	name := *params.RepositoryName
	if f.Repositories[name] != nil {
		return nil, errors.Errorf("name %s already exists", name)
	}

	now := time.Now()
	uri := "myawssession.dkr.ecr.myregion.amazonaws.com"
	id := uri
	repo := &types.Repository{
		CreatedAt:      &now,
		RegistryId:     &id,
		RepositoryArn:  nil,
		RepositoryName: &name,
		RepositoryUri:  &uri,
	}
	f.Repositories[name] = repo

	return &ecr.CreateRepositoryOutput{
		Repository:     repo,
		ResultMetadata: middleware.Metadata{},
	}, nil
}

// NewFakeECR creates a new fake ECR
func NewFakeECR() *FakeECR {
	return &FakeECR{
		Repositories: map[string]*types.Repository{},
	}
}
