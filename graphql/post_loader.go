package graphql

import (
  "context"
  "database/sql"
  "net/http"
  "sync"
  "time"
)

type PostLoader struct {
  criteria map[string]*Pagination
  data     map[string][]Post
  db       *sql.DB
  mutex    sync.Mutex
}

func NewPostLoader(db *sql.DB) *PostLoader {
  return &PostLoader{
    criteria: make(map[string]*Pagination),
    data:     make(map[string][]Post),
    db:       db,
    mutex:    sync.Mutex{},
  }
}

func (p *PostLoader) Query(ctx context.Context, userID string, pagination *Pagination) func() ([]Post, error) {
  p.mutex.Lock()
  if posts, ok := p.data[userID]; ok {
    return func() ([]Post, error) {
      return posts, nil
    }
  }
  p.criteria[userID] = pagination
  p.mutex.Unlock()

  finish := make(chan struct{})
  go func() {
    <-time.After(5 * time.Millisecond)
    finish <- struct{}{}
  }()

  return func() ([]Post, error) {
    <-finish
    err := p.load()
    if err != nil {
      return nil, err
    }
    return nil, nil
  }
}

func (p *PostLoader) load() error {
  /*
	for userID, criteria := range p.criteria {
	}*/
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
