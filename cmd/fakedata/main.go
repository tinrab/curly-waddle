package main

import (
	"database/sql"
	"log"
	"math/rand"
	"time"

	"github.com/icrowley/fake"
	"github.com/lib/pq"
	"github.com/segmentio/ksuid"
)

func main() {
	db, err := sql.Open("postgres", "postgres://blog:123456@127.0.0.1:5432/blog?sslmode=disable")
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

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// Insert users
	insertUser, err := tx.Prepare(pq.CopyIn("users", "id", "name"))
	if err != nil {
		log.Fatal(err)
	}
	var userIDs []string
	const numUsers = 10000
	for i := 1; i <= numUsers; i++ {
		id := ksuid.New().String()
		userIDs = append(userIDs, id)
		_, err := insertUser.Exec(id, fake.FullName())
		if err != nil {
			log.Fatal(err)
		}
	}
	_, err = insertUser.Exec()
	if err != nil {
		log.Fatal(err)
	}

	// Insert posts
	insertPost, err := tx.Prepare(pq.CopyIn("posts", "id", "user_id", "created_at", "body"))
	if err != nil {
		log.Fatal(err)
	}
	for i := 1; i <= 100; i++ {
		numPosts := int(rand.Int31n(10000) + 2)
		for j := 1; j <= numPosts; j++ {
			offset := time.Duration(rand.Int31n(1000000)) * time.Minute
			_, err := insertPost.Exec(
				ksuid.New().String(),
				userIDs[rand.Int31n(int32(len(userIDs)))],
				time.Now().Add(offset).UTC(),
				fake.Sentence(),
			)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	_, err = insertPost.Exec()
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}
