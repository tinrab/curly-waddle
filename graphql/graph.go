//go:generate gorunpkg github.com/vektah/gqlgen

package graphql

import (
  "database/sql"
)

type GraphQLServer struct {
  db *sql.DB
}

func NewGraphQLServer(db *sql.DB) *GraphQLServer {
  return &GraphQLServer{
    db: db,
  }
}

func (s *GraphQLServer) Mutation() MutationResolver {
  return &mutationResolver{
    server: s,
  }
}

func (s *GraphQLServer) Query() QueryResolver {
  return &queryResolver{
    server: s,
  }
}

func (s *GraphQLServer) User() UserResolver {
  return &userResolver{
    server: s,
  }
}
