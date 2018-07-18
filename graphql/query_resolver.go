package graphql

import (
	"context"
)

type queryResolver struct {
	server *GraphQLServer
}

func (r *queryResolver) Users(ctx context.Context, pagination *Pagination) ([]User, error) {
	if pagination == nil {
		pagination = &Pagination{
			Skip: 0,
			Take: 100,
		}
	}

	rows, err := r.server.db.QueryContext(
		ctx,
		"SELECT id, name FROM users OFFSET $1 LIMIT $2",
		pagination.Skip,
		pagination.Take,
	)
	if err != nil {
		return nil, err
	}

	user := &User{}
	var users []User

	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Name)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}

	return users, nil
}

func (r *queryResolver) Posts(ctx context.Context, pagination *Pagination) ([]Post, error) {
	if pagination == nil {
		pagination = &Pagination{
			Skip: 0,
			Take: 100,
		}
	}

	rows, err := r.server.db.QueryContext(
		ctx,
		"SELECT id, user_id, created_at, body FROM posts OFFSET $1 LIMIT $2",
		pagination.Skip,
		pagination.Take,
	)
	if err != nil {
		return nil, err
	}

	post := &Post{
		User: User{},
	}
	var posts []Post

	for rows.Next() {
		err = rows.Scan(&post.ID, &post.User.ID, &post.CreatedAt, &post.Body)
		if err != nil {
			return nil, err
		}
		posts = append(posts, *post)
	}

	return posts, nil
}
