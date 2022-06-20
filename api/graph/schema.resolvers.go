package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"linkshare_api/database"
	"linkshare_api/graph/generated"
	"linkshare_api/graph/model"
	"linkshare_api/utils"
)

func (r *mutationResolver) CreatePage(ctx context.Context, input model.NewPage) (*model.Page, error) {
	db, err := database.NewLinkShareDB(ctx)
	if err != nil {
		utils.LogError(err.Error())
		return nil, err
	}
	return db.CreateNewPage(ctx, input.URL, input.User)
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateUser(ctx context.Context, input model.UpdateUser) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdatePage(ctx context.Context, input model.UpdatePage) (*model.Page, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) User(ctx context.Context, username string) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Page(ctx context.Context, url string) (*model.Page, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) Pages(ctx context.Context, obj *model.User) ([]*model.Page, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
