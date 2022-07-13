package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"linkshare_api/contextual"
	"linkshare_api/database"
	"linkshare_api/graph/generated"
	"linkshare_api/graph/model"
	"linkshare_api/utils"
)

func (r *mutationResolver) CreatePage(ctx context.Context, url string) (*model.Page, error) {
	user := contextual.UserForContext(ctx)
	if user == nil {
		return nil, errors.New("must login to create page")
	}
	db, err := database.NewLinkShareDB(ctx)
	if err != nil {
		utils.LogError(err.Error())
		return nil, err
	}
	return db.CreatePage(ctx, url, user.ID, db.Pages.InsertOne, db.Users.UpdateByID)
}

func (r *mutationResolver) UpdateUser(ctx context.Context, input model.UpdateUser) (*model.User, error) {
	user := contextual.UserForContext(ctx)
	if user == nil {
		return nil, errors.New("must login to create page")
	}
	db, err := database.NewLinkShareDB(ctx)
	if err != nil {
		utils.LogError(err.Error())
		return nil, err
	}
	if input.Email != nil {
		user.Email = input.Email
	}
	if input.FirstName != nil {
		user.FirstName = input.FirstName
	}
	if input.LastName != nil {
		user.LastName = input.LastName
	}
	if input.Username != nil {
		user.Username = input.LastName
	}
	err = user.Update(ctx, db.Users.UpdateByID)

	return user, err
}

func (r *mutationResolver) UpdatePage(ctx context.Context, input model.UpdatePage) (*model.Page, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeletePage(ctx context.Context, url string) (*bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) User(ctx context.Context, username string) (*model.User, error) {
	db, err := database.NewLinkShareDB(ctx)
	if err != nil {
		utils.LogError(err.Error())
		return nil, err
	}
	user := &model.User{
		Username: &username,
	}
	err = user.LoadByUsername(ctx, db.Users.FindOne)
	return user, err
}

func (r *queryResolver) Page(ctx context.Context, url string) (*model.Page, error) {
	db, err := database.NewLinkShareDB(ctx)
	if err != nil {
		utils.LogError(err.Error())
		return nil, err
	}
	page := &model.Page{
		URL: url,
	}
	page.LoadByURL(ctx, db.Pages.FindOne)
	return page, err
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
