package main

import (
  "database/sql"
  "log"
  "math/rand"
  "strings"
  "time"

  _ "github.com/go-sql-driver/mysql"
  "github.com/icrowley/fake"
  "github.com/segmentio/ksuid"
)

func main() {
  db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/blog")
  if err != nil {
    log.Fatal(err)
  }
  err = db.Ping()
  if err != nil {
    log.Fatal(err)
  }

  // Drop
  _, err = db.Exec("DELETE FROM users")
  _, err = db.Exec("DELETE FROM posts")
  if err != nil {
    log.Fatal(err)
  }

  // Insert users
  var args []interface{}
  q := strings.Builder{}
  var userIDs []string

  q.WriteString("INSERT INTO users(id, name) VALUES")

  const numUsers = 10000
  for i := 1; i <= numUsers; i++ {
    q.WriteString("(?, ?)")
    if i <= numUsers-1 {
      q.WriteRune(',')
    }

    id := ksuid.New().String()
    userIDs = append(userIDs, id)
    args = append(args, id)
    args = append(args, fake.FullName())
  }
  _, err = db.Exec(q.String(), args...)
  if err != nil {
    log.Fatal(err)
  }

  // Insert posts
  q = strings.Builder{}
  for i := 1; i <= 100; i++ {
    args = []interface{}{}
    q.Reset()
    q.WriteString("INSERT INTO posts(id, user_id, body, created_at) VALUES")
    numPosts := int(rand.Int31n(10000) + 2)
    for j := 1; j <= numPosts; j++ {
      q.WriteString("(?, ?, ?, ?)")
      if j <= numPosts-1 {
        q.WriteRune(',')
      }
      args = append(args, ksuid.New().String())
      args = append(args, userIDs[rand.Int31n(int32(len(userIDs)))])
      args = append(args, fake.Sentence())
      offset := time.Duration(rand.Int31n(1000000)) * time.Minute
      args = append(args, time.Now().Add(offset).UTC())
    }

    _, err = db.Exec(q.String(), args...)
    if err != nil {
      log.Fatal(err)
    }
  }
}
