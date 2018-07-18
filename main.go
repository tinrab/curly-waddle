package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/tinrab/curly-waddle/graphql"
	"github.com/vektah/gqlgen/handler"
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

	s := graphql.NewGraphQLServer(db)
	http.Handle("/graphql", graphql.PostLoaderMiddleware(db, handler.GraphQL(graphql.NewExecutableSchema(s))))

	log.Println("Running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
