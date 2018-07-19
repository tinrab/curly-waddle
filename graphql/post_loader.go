package graphql

import (
	"context"
	"database/sql"
	"sync"

	"github.com/lib/pq"
)

type PostLoader struct {
	pagination Pagination
	data       map[string][]Post
	db         *sql.DB
	mutex      sync.Mutex
	loaded     bool
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

func (p *PostLoader) Enqueue(userIDs []string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, userID := range userIDs {
		p.data[userID] = []Post{}
	}
}

func (p *PostLoader) Query(ctx context.Context, userID string, pagination *Pagination) ([]Post, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.loaded {
		return p.data[userID], nil
	}

	if len(p.data) == 0 {
		return nil, nil
	}

	if pagination != nil {
		p.pagination = *pagination
	}

	var userIDs []string
	for userID := range p.data {
		userIDs = append(userIDs, userID)
	}
	err := p.load(userIDs)
	if err != nil {
		return nil, err
	}
	p.loaded = true

	userLoader := ctx.Value(userLoaderKey{}).(*UserLoader)
	if userLoader != nil {
		userLoader.Enqueue(userIDs)
	}

	return p.data[userID], nil
}

func (p *PostLoader) load(userIDs []string) error {
	rows, err := p.db.Query(
		"SELECT id, user_id, created_at, body FROM read_user_posts ($1, $2, $3)",
		pq.Array(userIDs),
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

		p.data[userID] = append(p.data[userID], Post{
			ID:        post.ID,
			CreatedAt: post.CreatedAt,
			Body:      post.Body,
			User: &User{
				ID: userID,
			},
		})
	}

	return nil
}
