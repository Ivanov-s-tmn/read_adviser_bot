package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"read-adviser-bot/lib/e"
	"read-adviser-bot/storage"
	"time"
)

type Storage struct {
	basePath string
}

const defaultPermission = 0774

func New(basePath string) Storage {
	return Storage{
		basePath: basePath,
	}
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() { err = e.WrapIfErr("can't save page", err) }()

	filePath := filepath.Join(s.basePath, page.Username)

	if err := os.MkdirAll(filePath, defaultPermission); err != nil {
		return err
	}

	fileName, err := fileName(page)
	if err != nil {
		return err
	}

	filePath = filepath.Join(filePath, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil
}

func (s Storage) PickRandom(username string) (page *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("can't pick random page", err) }()

	path := filepath.Join(s.basePath, username)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	rand.NewSource(time.Now().UnixNano())
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return e.Wrap("can't remove file", err)
	}

	path := filepath.Join(s.basePath, p.Username, fileName)

	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("can't remove file %s", path)

		return e.Wrap(msg, err)
	}
	return nil
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("can't check if file exists", err)
	}

	path := filepath.Join(s.basePath, p.Username, fileName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file %s exists", path)

		return false, e.Wrap(msg, err)
	}
	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("can't decode page", err)
	}
	defer func() { _ = f.Close() }()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, err
	}

	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
