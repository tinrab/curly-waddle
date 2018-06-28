//go:generate gqlgen -schema ./schema.graphql
package graphql

import (
	"context"
	"database/sql"
)

type graphQLServer struct {
	db *sql.DB
}

func NewGraphQLServer(db *sql.DB) *graphQLServer {
	return &graphQLServer{
		db: db,
	}
}

func (*graphQLServer) Mutation_createUser(ctx context.Context, input CreateUserInput) (*User, error) {
	return nil, nil
}

func (*graphQLServer) Mutation_createPost(ctx context.Context, input CreatePostInput) (*Post, error) {
	return nil, nil
}

func (*graphQLServer) Query_users(ctx context.Context, skip *int, take *int) ([]User, error) {
	return nil, nil
}

func (*graphQLServer) Query_posts(ctx context.Context, skip *int, take *int) ([]Post, error) {
	return nil, nil
}
