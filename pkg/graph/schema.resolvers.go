package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"math"

	"github.com/b-open/jobbuzz/pkg/graph/generated"
	"github.com/b-open/jobbuzz/pkg/graph/graphmodel"
	"github.com/jinzhu/copier"
	werrors "github.com/pkg/errors"
)

func (r *mutationResolver) RegisterAccount(ctx context.Context, input graphmodel.NewUserInput) (graphmodel.NewUser, error) {
	token, err := r.Service.CreateUser(input.Email, input.Password)
	if err != nil {
		return nil, werrors.Wrapf(err, "Error in CreateUser")
	}

	output := graphmodel.LoginResult{
		AccessToken: token,
	}

	return output, nil
}

func (r *queryResolver) Jobs(ctx context.Context, search *graphmodel.StringFilterInput, pagination graphmodel.PaginationInput) (*graphmodel.JobOutput, error) {
	jobs, totalCount, err := r.Service.GetJobs(pagination)
	if err != nil {
		return nil, werrors.Wrapf(err, "Error in GetJobs")
	}

	var graphqlJobs []*graphmodel.Job
	err = copier.Copy(&graphqlJobs, &jobs)
	if err != nil {
		return nil, werrors.Wrapf(err, "Error copying structs")
	}

	// TODO: move pagination logic to a func
	output := &graphmodel.JobOutput{
		From:        pagination.Offset,
		To:          pagination.Offset + len(graphqlJobs),
		PerPage:     pagination.Limit,
		CurrentPage: int(math.Ceil(float64(pagination.Offset) / float64(pagination.Limit))),
		TotalPage:   int(math.Ceil(float64(totalCount) / float64(pagination.Limit))),
		Total:       int(totalCount),
		Data:        graphqlJobs,
	}

	return output, nil
}

func (r *queryResolver) Companies(ctx context.Context, search *graphmodel.StringFilterInput, pagination graphmodel.PaginationInput) (*graphmodel.CompanyOutput, error) {
	companies, totalCount, err := r.Service.GetCompanies(pagination)
	if err != nil {
		return nil, werrors.Wrapf(err, "Error in GetCompanies")
	}

	var graphqlCompanies []*graphmodel.Company
	err = copier.Copy(&graphqlCompanies, &companies)
	if err != nil {
		return nil, werrors.Wrapf(err, "Error copying structs")
	}

	output := &graphmodel.CompanyOutput{
		From:        pagination.Offset,
		To:          pagination.Offset + len(graphqlCompanies),
		PerPage:     pagination.Limit,
		CurrentPage: int(math.Ceil(float64(pagination.Offset) / float64(pagination.Limit))),
		TotalPage:   int(math.Ceil(float64(totalCount) / float64(pagination.Limit))),
		Total:       int(totalCount),
		Data:        graphqlCompanies,
	}

	return output, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
