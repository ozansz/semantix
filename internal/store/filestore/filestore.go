package filestore

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/ozansz/semantix/internal/parser"
	"github.com/ozansz/semantix/internal/store"

	"github.com/zhenjl/cityhash"
)

const (
	defaultStoreFileMode = 0666
	defaultFlushInterval = 5 * time.Second
)

type FileStore struct {
	fp            *os.File
	debug         bool
	path          string
	persistent    bool
	store         sync.Map
	idBuffer      sync.Map
	storeSyncDone chan struct{}
}

type FileStoreOption func(*FileStore)

func WithPersistentFile(path string) FileStoreOption {
	return func(fs *FileStore) {
		fs.path = path
		fs.persistent = true
	}
}

func WithDebug() FileStoreOption {
	return func(fs *FileStore) {
		fs.debug = true
	}
}

func New(opts ...FileStoreOption) (*FileStore, error) {
	fs := &FileStore{
		store:         sync.Map{},
		idBuffer:      sync.Map{},
		storeSyncDone: make(chan struct{}),
		debug:         false,
	}
	for _, o := range opts {
		o(fs)
	}

	var err error
	if fs.persistent {
		fs.fp, err = os.Create(fs.path)
		if err != nil {
			return nil, fmt.Errorf("os.OpenFile: %v", err)
		}
		go fs.syncWorker(fs.storeSyncDone)
	}

	return fs, nil
}

func (fs *FileStore) syncWorker(done <-chan struct{}) {
	t := time.NewTicker(defaultFlushInterval)
	for {
		select {
		case <-done:
			fs.flush()
			if fs.debug {
				fmt.Printf("FileStore worker stopped!")
			}
			return
		case <-t.C:
			fs.flush()
		}
	}
}

func (fs *FileStore) flush() {
	if fs.persistent {
		errs := []error{}

		fs.idBuffer.Range(func(key, value any) bool {
			id := key.(uint32)
			t, ok := fs.store.Load(id)
			if !ok {
				errs = append(errs, fmt.Errorf("triple with id %d not found in store", id))
				return true
			}
			enc := tripleEncode(id, t.(*parser.Fact))
			if _, err := fs.fp.Write(enc); err != nil {
				errs = append(errs, err)
			}
			if _, err := fs.fp.Write([]byte("\n")); err != nil {
				errs = append(errs, err)
			}
			return true
		})
		fs.idBuffer = sync.Map{}
		if err := fs.fp.Sync(); err != nil {
			errs = append(errs, err)
		}
		if len(errs) > 0 {
			e := ""
			for _, err := range errs {
				e += fmt.Sprintf("- %v\n", err)
			}
			if fs.debug {
				fmt.Printf("!!! Error writing to file at filestore.FileStore.flush:\n%s\n", e)
			} else {
				log.Panicf("!!! Error writing to file at filestore.FileStore.flush:\n%s\n", e)
			}
		}
	}
}

func (fs *FileStore) Close() error {
	if fs.persistent {
		close(fs.storeSyncDone)
		return fs.fp.Close()
	}
	return nil
}

func (fs *FileStore) Add(t *parser.Fact) error {
	h := tripleHash(t)
	fs.store.Store(h, t.Copy())
	fs.idBuffer.Store(h, true)
	return nil
}

func (fs *FileStore) Get(q *store.Query) (map[uint32]*parser.Fact, error) {
	trs := map[uint32]*parser.Fact{}

	fs.store.Range(func(key, value any) bool {
		t := value.(*parser.Fact)
		if q.Matches(t) {
			trs[key.(uint32)] = t.Copy()
		}
		return true
	})

	return trs, nil
}

func (fs *FileStore) Sync() error {
	fs.flush()
	return nil
}

func tripleHash(t *parser.Fact) uint32 {
	var s string
	if t.Subject != nil {
		if t.Object != nil {
			s = fmt.Sprintf("s:%q|p:%q|o:%d:%q", *t.Subject, t.Predicate, t.Object.Kind(), t.Object.String())
		} else {
			s = fmt.Sprintf("s:%q|p:%q|o:(%s)", *t.Subject, t.Predicate, t.ObjectFact.Pretty())
		}
	} else {
		if t.Object != nil {
			s = fmt.Sprintf("s:(%s)|p:%q|o:%d:%q", t.SubjectFact.Pretty(), t.Predicate, t.Object.Kind(), t.Object.String())
		} else {
			s = fmt.Sprintf("s:(%s)|p:%q|o:(%s)", t.SubjectFact.Pretty(), t.Predicate, t.ObjectFact.Pretty())
		}
	}
	b := []byte(s)
	return cityhash.CityHash32(b, uint32(len(b)))
}

func tripleEncode(id uint32, t *parser.Fact) []byte {
	var s string
	if t.Subject != nil {
		if t.Object != nil {
			s = fmt.Sprintf("(%d, %s, %s, %d, %s)", id, *t.Subject, t.Predicate, t.Object.Kind(), t.Object.String())
		} else {
			s = fmt.Sprintf("(%d, %s, %s, (%s))", id, *t.Subject, t.Predicate, t.ObjectFact.Pretty())
		}
	} else {
		if t.Object != nil {
			s = fmt.Sprintf("(%d, (%s), %s, %d, %s)", id, t.SubjectFact.Pretty(), t.Predicate, t.Object.Kind(), t.Object.String())
		} else {
			s = fmt.Sprintf("(%d, (%s), %s, (%s))", id, t.SubjectFact.Pretty(), t.Predicate, t.ObjectFact.Pretty())
		}
	}
	return []byte(s)
}
