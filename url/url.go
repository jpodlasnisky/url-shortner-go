package url

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"encore.dev/storage/sqldb"
)

type URL struct {
	ID  string // short-form URL id
	URL string // long-form URL
}

type ShortenParams struct {
	URL string // long-form URL to shorten
}

// Shorten shortens a URL
//
//encore:api public	method=POST path=/url
func Shorten(ctx context.Context, p *ShortenParams) (*URL, error) {
	id, err := generateID()
	if err != nil {
		return nil, err
	} else if err := insert(ctx, id, p.URL); err != nil {
		return nil, err
	}
	return &URL{ID: id, URL: p.URL}, nil
}

// generateID generates a random short ID.
func generateID() (string, error) {
	var data [6]byte // 6 bytes = 48 bits -> entropy
	if _, err := rand.Read(data[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data[:]), nil
}

// insert inserts a URL into the database.
func insert(ctx context.Context, id, url string) error {
	_, err := sqldb.Exec(ctx, `
		INSERT INTO url (id, original_url)
		VALUES ($1, $2)
	`, id, url)
	return err
}

// Get retrieves the original URL for the id.
//
//encore:api public method=GET path=/url/:id
func Get(ctx context.Context, id string) (*URL, error) {
	u := &URL{ID: id}
	err := sqldb.QueryRow(ctx, `
		SELECT original_url FROM url
		WHERE id = $1
	`, id).Scan(&u.URL)
	return u, err
}
