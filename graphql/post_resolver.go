package graphql

import (
	"context"
)

type postResolver struct {
	server *GraphQLServer
}

func (r *postResolver) User(ctx context.Context, obj *Post) (*User, error) {
	userLoader := ctx.Value(userLoaderKey{}).(*UserLoader)
	return userLoader.Query(ctx, obj.User.ID)
}
