package graphql

import (
	"context"
	"database/sql"
	"net/http"
	"sync"
	"time"

	"github.com/lib/pq"
)

type PostLoader struct {
	pagination Pagination
	userIDs    []string
	data       map[string][]Post
	db         *sql.DB
	mutex      sync.Mutex
}

func NewPostLoader(db *sql.DB) *PostLoader {
	return &PostLoader{
		pagination: Pagination{
			Skip: 0,
			Take: 10,
		},
		data:  make(map[string][]Post),
		db:    db,
		mutex: sync.Mutex{},
	}
}

func (p *PostLoader) Query(ctx context.Context, userID string, pagination *Pagination) func() ([]Post, error) {
	p.mutex.Lock()
	if posts, ok := p.data[userID]; ok {
		return func() ([]Post, error) {
			return posts, nil
		}
	}
	if pagination != nil {
		p.pagination = *pagination
	}
	p.userIDs = append(p.userIDs, userID)
	p.mutex.Unlock()

	finish := make(chan struct{})
	go func() {
		<-time.After(5 * time.Millisecond)
		finish <- struct{}{}
	}()

	return func() ([]Post, error) {
		<-finish
		if posts, ok := p.data[userID]; ok {
			return posts, nil
		}
		err := p.load()
		if err != nil {
			return nil, err
		}
		return p.data[userID], nil
	}
}

func (p *PostLoader) load() error {
	if len(p.userIDs) == 0 {
		return nil
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	rows, err := p.db.Query(
		"SELECT id, user_id, created_at, body FROM read_user_posts ($1, $2, $3)",
		pq.Array(p.userIDs),
		p.pagination.Skip,
		p.pagination.Take,
	)
	if err != nil {
		return err
	}

	post := Post{}
	var userID string
	for rows.Next() {
		err = rows.Scan(&post.ID, &userID, &post.CreatedAt, &post.Body)
		if err != nil {
			return err
		}

		if _, ok := p.data[userID]; !ok {
			p.data[userID] = []Post{}
		}
		p.data[userID] = append(p.data[userID], Post{
			ID:        post.ID,
			CreatedAt: post.CreatedAt,
			Body:      post.Body,
			User: User{
				ID: userID,
			},
		})
	}

	p.userIDs = []string{}

	return nil
}

type postLoaderKey struct{}

func PostLoaderMiddleware(db *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postLoader := NewPostLoader(db)

		ctx := context.WithValue(r.Context(), postLoaderKey{}, postLoader)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
