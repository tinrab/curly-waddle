package graphql

import (
	"context"
	"database/sql"
	"net/http"
)

type postLoaderKey struct{}

type userLoaderKey struct{}

func PostLoaderMiddleware(db *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), postLoaderKey{}, NewPostLoader(db))
		ctx = context.WithValue(ctx, userLoaderKey{}, NewUserLoader(db))
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
