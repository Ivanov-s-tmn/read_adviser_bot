package storage

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"read-adviser-bot/lib/e"
	"time"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(username string) (*Page, error)
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
}

var ErrNoSavedPages = errors.New("no saved pages")

type Page struct {
	URL      string
	Username string
	Created  time.Time
}

func (p Page) Hash() (hash string, err error) {
	defer func() { err = e.WrapIfErr("can't calculate hash", err) }()
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", err
	}

	if _, err := io.WriteString(h, p.Username); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
