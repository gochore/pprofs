package pprofs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	EnvPprofDir = "PPROF_DIR"
)

const (
	pprofFileSuffix = ".pb.gz"
)

type Storage interface {
	WriteCloser(name string, t time.Time) (io.WriteCloser, error)
}

type FileStorage struct {
	prefix    string
	dir       string
	ttl       time.Duration
	cleaning  bool
	cleaningM *sync.Mutex
}

func NewFileStorage(prefix string, dir string, ttl time.Duration) *FileStorage {
	return &FileStorage{
		prefix:    prefix,
		dir:       dir,
		ttl:       ttl,
		cleaning:  false,
		cleaningM: &sync.Mutex{},
	}
}

func NewFileStorageFromEnv() *FileStorage {
	prefix := os.Args[0]

	dir := "/tmp/pprofs"
	if v := os.Getenv(EnvPprofDir); v != "" {
		dir = v
	}

	return NewFileStorage(prefix, dir, 24*time.Hour)
}

func (s *FileStorage) WriteCloser(name string, t time.Time) (io.WriteCloser, error) {
	file := fmt.Sprintf("%s-%s.%s"+pprofFileSuffix, s.prefix, t.Format("20060102T150405"), name)
	path := filepath.Join(s.dir, file)
	if err := os.MkdirAll(s.dir, 0755); err != nil {
		return nil, fmt.Errorf("mkdir %q: %w", s.dir, err)
	}

	return s.newWriteCloseCleaner(path), nil
}

func (s *FileStorage) clean() {
	cleaning := false

	s.cleaningM.Lock()
	cleaning = s.cleaning
	if !s.cleaning {
		s.cleaning = true
	}
	s.cleaningM.Unlock()

	if cleaning {
		return
	}

	files, err := os.ReadDir(s.dir)
	if err != nil {
		return
	}
	for _, v := range files {
		if v.IsDir() {
			continue
		}
		if name := v.Name(); strings.HasSuffix(name, pprofFileSuffix) {
			if info, err := v.Info(); err == nil {
				if time.Since(info.ModTime()) > s.ttl {
					_ = os.Remove(filepath.Join(s.dir, name))
				}
			}
		}
	}

	s.cleaningM.Lock()
	s.cleaning = false
	s.cleaningM.Unlock()
}

func (s *FileStorage) newWriteCloseCleaner(path string) *writeCloseCleaner {
	return &writeCloseCleaner{
		path:  path,
		clean: s.clean,
	}
}

type writeCloseCleaner struct {
	path  string
	file  *os.File
	clean func()
}

func (c *writeCloseCleaner) Write(p []byte) (n int, err error) {
	if c.file == nil {
		f, err := os.Create(c.path)
		if err != nil {
			return 0, fmt.Errorf("create %q: %w", c.path, err)
		}
		c.file = f
	}
	return c.file.Write(p)
}

func (c *writeCloseCleaner) Close() error {
	var err error
	if c.file != nil {
		err = c.file.Close()
	}
	c.clean()
	return err
}
