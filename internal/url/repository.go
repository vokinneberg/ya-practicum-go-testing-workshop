package url

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	upsertQuery = `
		INSERT INTO short_urls (id, original_url, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (id)
		DO UPDATE SET
			original_url = EXCLUDED.original_url,
			created_at = EXCLUDED.created_at
		RETURNING id, original_url, created_at;
	`
	selectByIDQuery = `
		SELECT id, original_url, created_at FROM short_urls WHERE id = $1;
	`
	selectAllQuery = `
		SELECT id, original_url, created_at FROM short_urls;
	`
)

// Repository is the repository for the short URLs
type Repository struct {
	dbpool    *pgxpool.Pool
	cache     map[string]*ShortURL
	syncWrite bool
}

// NewRepository creates a new repository
func NewRepository(dbpool *pgxpool.Pool, syncWrite bool) *Repository {
	return &Repository{
		dbpool:    dbpool,
		cache:     make(map[string]*ShortURL),
		syncWrite: syncWrite,
	}
}

// CreateURL creates a short URL
func (r *Repository) CreateURL(ctx context.Context, url *ShortURL) (*ShortURL, error) {
	// Check if the URL is in the cache
	if shortURL, ok := r.cache[url.ID]; ok {
		return shortURL, nil
	}

	// If syncWrite is true, write to the database synchronously
	if r.syncWrite {
		createdAt := time.Now()
		err := r.dbpool.QueryRow(ctx, upsertQuery, url.ID, url.OriginalURL, createdAt).Scan(&url.ID, &url.OriginalURL, &url.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to insert short url: %w", err)
		}
	}

	// Add the URL to the cache
	r.cache[url.ID] = url
	return url, nil
}

// GetURL returns a short URL by its ID
func (r *Repository) GetURL(ctx context.Context, id string) (*ShortURL, error) {
	// Check if the URL is in the cache
	if url, ok := r.cache[id]; ok {
		return url, nil
	}

	// If not, get it from the database
	var url ShortURL
	err := r.dbpool.QueryRow(ctx, selectByIDQuery, id).Scan(&url.ID, &url.OriginalURL, &url.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get short url: %w", err)
	}

	// Add the URL to the cache
	r.cache[id] = &url
	return &url, nil
}

func (r *Repository) Restore(ctx context.Context) error {
	rows, err := r.dbpool.Query(ctx, selectAllQuery)
	if err != nil {
		return fmt.Errorf("failed to get short urls: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var url ShortURL
		if err := rows.Scan(&url.ID, &url.OriginalURL, &url.CreatedAt); err != nil {
			return fmt.Errorf("failed to scan short url: %w", err)
		}
		r.cache[url.ID] = &url
	}

	return nil
}

// Flush dumps the cached data to the database
func (r *Repository) Flush(ctx context.Context) error {
	for _, url := range r.cache {
		createdAt := time.Now()
		err := r.dbpool.QueryRow(ctx, upsertQuery, url.ID, url.OriginalURL, createdAt).Scan(&url.ID, &url.OriginalURL, &url.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to insert short url: %w", err)
		}
	}
	return nil
}
