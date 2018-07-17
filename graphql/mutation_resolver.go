package graphql

import (
  "context"
  "time"

  "github.com/segmentio/ksuid"
)

type mutationResolver struct {
  server *GraphQLServer
}

func (r *mutationResolver) CreateUser(ctx context.Context, input CreateUserInput) (*User, error) {
  user := &User{
    ID:   ksuid.New().String(),
    Name: input.Name,
  }
  _, err := r.server.db.ExecContext(ctx, "INSERT INTO `users`(`id`, `name`) VALUES(?, ?)", user.ID, user.Name)
  if err != nil {
    return nil, err
  }
  return user, nil
}

func (r *mutationResolver) CreatePost(ctx context.Context, input CreatePostInput) (*Post, error) {
  user := &User{}
  err := r.server.db.
    QueryRowContext(ctx, "SELECT `id`, `name` FROM `users` WHERE `id`=?", input.UserId).
    Scan(&user.ID, &user.Name)
  if err != nil {
    return nil, err
  }

  post := &Post{
    ID:        ksuid.New().String(),
    CreatedAt: time.Now().UTC(),
    Body:      input.Body,
    User:      *user,
  }
  _, err = r.server.db.ExecContext(
    ctx,
    "INSERT INTO `posts`(`id`, `user_id`, `created_at`, `body`) VALUES(?, ?, ?, ?)",
    post.ID,
    post.User.ID,
    post.CreatedAt,
    post.Body,
  )
  if err != nil {
    return nil, err
  }

  return post, nil
}
