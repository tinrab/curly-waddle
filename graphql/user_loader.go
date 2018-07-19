package graphql

import (
	"context"
	"database/sql"
	"sync"

	"github.com/lib/pq"
)

type UserLoader struct {
	data   map[string]*User
	loaded bool
	db     *sql.DB
	mutex  sync.Mutex
}

func NewUserLoader(db *sql.DB) *UserLoader {
	return &UserLoader{
		data:  make(map[string]*User),
		db:    db,
		mutex: sync.Mutex{},
	}
}

func (u *UserLoader) Enqueue(userIDs []string) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	for _, userID := range userIDs {
		u.data[userID] = &User{}
	}
}

func (u *UserLoader) Query(ctx context.Context, userID string) (*User, error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if u.loaded {
		return u.data[userID], nil
	}

	if len(u.data) == 0 {
		return nil, nil
	}

	var userIDs []string
	for userID := range u.data {
		userIDs = append(userIDs, userID)
	}
	err := u.load(userIDs)
	if err != nil {
		return nil, err
	}
	u.loaded = true

	postLoader := ctx.Value(postLoaderKey{}).(*PostLoader)
	if postLoader != nil {
		postLoader.Enqueue(userIDs)
	}

	return u.data[userID], nil
}

func (u *UserLoader) load(userIDs []string) error {
	rows, err := u.db.Query(
		"SELECT id, name FROM users WHERE id = ANY($1::CHAR(27)[])",
		pq.Array(userIDs),
	)
	if err != nil {
		return err
	}

	for rows.Next() {
		user := &User{}
		err = rows.Scan(&user.ID, &user.Name)
		if err != nil {
			return err
		}
		u.data[user.ID] = user
	}

	return nil
}
