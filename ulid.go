// Package safeulid provides concurrent safe functions for generating monotonic ulids
package safeulid

import (
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

// EntropyReader is a factory for creating new sources of random entropy.
type EntropyReader struct {
	New func() io.Reader
}

// DefaultEntropyReader is the default EntropyReader factory.
var DefaultEntropyReader = &EntropyReader{
	New: func() io.Reader {
		return rand.New(rand.NewSource(time.Now().UnixNano()))
	},
}

var ids = NewDefaultFactory()

// IDFactory represents a concurrently safe way to generate new ids
type IDFactory struct {
	pool sync.Pool
}

// NewFactory returns a new IDFactory which can be concurrently accessed to to generate monotonic ulids.
// If entropy is Nil DefaultEntropySource is used.
// Use only if you have the need to generate ids from multiple different entropy sources or else use New() or MustNew().
func NewFactory(entropy *EntropyReader) *IDFactory {
	if entropy == nil {
		panic("entropy must not be nil")
	}
	f := new(IDFactory)
	f.pool.New = func() interface{} {
		return ulid.Monotonic(entropy.New(), 0)
	}
	return f
}

// NewDefaultFactory returns a new IDFactory from a DefaultEntropyReader
func NewDefaultFactory() *IDFactory {
	// return NewFactory(newSafeRand())
	return NewFactory(DefaultEntropyReader)
}

// New returns a monotonic ID and an error.
func (i *IDFactory) New() (ID, error) {
	return i.new(clck{})

}

// MustNew like New but panics on error.
func (i *IDFactory) MustNew() ID {
	return i.mustNew(clck{})
}

func (i *IDFactory) mustNew(t clock) ID {
	id, err := i.new(t)
	if err != nil {
		panic(err)
	}
	return id

}

func (i *IDFactory) new(t clock) (ID, error) {
	me := i.pool.Get().(*ulid.MonotonicEntropy)
	id, err := ulid.New(ulid.Timestamp(t.Now()), me)
	i.pool.Put(me)
	return ID{id}, err
}

type ID struct {
	ulid.ULID
}

// MustNew returns a new monotonic ID in a concurrently safe way.
// Like New but panics on error.
func MustNew() ID {
	return ids.MustNew()
}

// New returns a new monotonic ID in a concurrently safe way.
// Internally this function uses the time pkg from the standard library.
// An error is returned only if timestamp returned by time.Now() is before the year 1678 or after the year 2262 or DefaultEntropyReader errors
func New() (ID, error) {
	return ids.New()
}

type clock interface {
	Now() time.Time
}

type clck struct{}

func (clck) Now() time.Time { return time.Now() }
