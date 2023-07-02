package filestore

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

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
	// subjectIndex
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
			enc := tripleEncode(id, t.(*store.Triple))
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

func (fs *FileStore) Add(t *store.Triple) error {
	h := tripleHash(t)
	fs.store.Store(h, t.Copy())
	fs.idBuffer.Store(h, true)
	return nil
}

func (fs *FileStore) Get(q *store.Query) (map[uint32]*store.Triple, error) {
	trs := map[uint32]*store.Triple{}

	fs.store.Range(func(key, value any) bool {
		t := value.(*store.Triple)
		if (q.SubjectFilter != nil && t.Subject != *q.SubjectFilter) ||
			(q.PredicateFilter != nil && t.Predicate != *q.PredicateFilter) ||
			(q.ObjectFilterString != nil && t.Object.Kind != store.ObjectKindFloat && *t.Object.StringValue != *q.ObjectFilterString) ||
			// (q.ObjectFilterInteger != nil && t.Object.Kind == store.ObjectKindInteger && *t.Object.IntegerValue != *q.ObjectFilterInteger) ||
			(q.ObjectFilterFloat != nil && t.Object.Kind == store.ObjectKindFloat && *t.Object.FloatValue != *q.ObjectFilterFloat) {
			// (q.ObjectFilterBool != nil && t.Object.Kind == store.ObjectKindBool && *t.Object.BoolValue != *q.ObjectFilterBool) {
			return true
		}
		trs[key.(uint32)] = t.Copy()
		return true
	})

	return trs, nil
}

func tripleHash(t *store.Triple) uint32 {
	s := fmt.Sprintf("s:%q|p:%q|o:%d:%q", t.Subject, t.Predicate, t.Object.Kind, t.Object.String())
	b := []byte(s)
	return cityhash.CityHash32(b, uint32(len(b)))
}

func tripleEncode(id uint32, t *store.Triple) []byte {
	s := fmt.Sprintf("(%d, %s, %s, %d, %s)", id, t.Subject,
		t.Predicate, t.Object.Kind, t.Object.String())
	return []byte(s)
}
