package graphql

import (
  "context"
)

type userResolver struct {
  server *GraphQLServer
}

func (r *userResolver) Posts(ctx context.Context, obj *User, pagination *Pagination) ([]Post, error) {
  return nil, nil
}
