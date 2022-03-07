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
	EnvPprofPrefix = "PPROF_PREFIX"
	EnvPprofDir    = "PPROF_DIR"
	EnvPprofTtl    = "PPROF_TTL"
)

const (
	pprofFileSuffix = ".pb.gz"
)

// Storage is an abstraction of storage output.
type Storage interface {
	// WriteCloser should return an io.WriteCloser according to the profile name and current time.
	WriteCloser(name string, t time.Time) (io.WriteCloser, error)
}

// FileStorage stores the profile results to file with automatic cleaning.
type FileStorage struct {
	prefix    string
	dir       string
	ttl       time.Duration
	cleaning  bool
	cleaningM *sync.Mutex
}

// NewFileStorage return a FileStorage,
// prefix indicates the prefix of the files,
// dir indicates where to save the files,
// ttl indicates the time-to-live of the files.
func NewFileStorage(prefix, dir string, ttl time.Duration) *FileStorage {
	return &FileStorage{
		prefix:    prefix,
		dir:       dir,
		ttl:       ttl,
		cleaning:  false,
		cleaningM: &sync.Mutex{},
	}
}

// NewFileStorageFromEnv return a FileStorage,
// it read arguments from environment variables.
func NewFileStorageFromEnv() *FileStorage {
	prefix := filepath.Base(os.Args[0])
	if v := os.Getenv(EnvPprofPrefix); v != "" {
		prefix = v
	}

	dir := filepath.Join(os.TempDir(), "pprofs")
	if v := os.Getenv(EnvPprofDir); v != "" {
		dir = v
	}

	ttl := 24 * time.Hour
	if v := os.Getenv(EnvPprofTtl); v != "" {
		if t, err := time.ParseDuration(v); err != nil {
			ttl = t
		}
	}

	return NewFileStorage(prefix, dir, ttl)
}

// WriteCloser implements Storage.
func (s *FileStorage) WriteCloser(name string, t time.Time) (io.WriteCloser, error) {
	if err := os.MkdirAll(s.dir, 0755); err != nil {
		return nil, fmt.Errorf("mkdir %q: %w", s.dir, err)
	}

	file := fmt.Sprintf("%s-%s.%s"+pprofFileSuffix, s.prefix, t.Format("20060102T150405"), name)
	path := filepath.Join(s.dir, file)
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
