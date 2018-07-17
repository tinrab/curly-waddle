package graphql

import (
  "time"
)

type Post struct {
  ID        string    `json:"id"`
  Body      string    `json:"body"`
  CreatedAt time.Time `json:"createdAt"`
  User      User      `json:"user"`
}

type User struct {
  ID    string `json:"id"`
  Name  string `json:"name"`
  Posts []Post `json:"posts"`
}
