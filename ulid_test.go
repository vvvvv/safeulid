package safeulid

import (
	"sync"
	"testing"
	"time"
)

type dummyClock struct {
	time.Time
}

func (d *dummyClock) Now() time.Time {
	return d.Time
}

var c = &dummyClock{time.Unix(1000000, 0)}

type constantReader struct {
	b byte
}

func (r constantReader) Read(b []byte) (int, error) {
	var i int
	for i = 0; i < len(b); i++ {
		b[i] = r.b
	}
	return i, nil
}

func TestNewFactory(t *testing.T) {
	t.Parallel()
	// r := constantReader{'x'}
	f := NewDefaultFactory()
	if f == nil {
		t.Fatalf("NewFactory should never return nil")
	}

	var id ID
	var err error
	id, err = f.new(c)
	if err != nil {
		t.Fatalf("NewFactory default should never return a factory that errors. err %v", err)
	}
	_ = id

	p := 2 // number of parallel coroutines
	number_of_ids := 5
	var mu sync.Mutex
	// all := make([]ID, p*number_of_ids)
	var all []ID
	var wg sync.WaitGroup
	for i := 0; i < p; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ids := make([]ID, number_of_ids)
			for j := 0; j < number_of_ids; j++ {
				// 0 * 100 + j 			= 	0
				//							1
				//							...
				//							99
				// 1 * 100 + j 			= 	100
				// 							101
				//							...
				// id = ulid.MustNew()
				// mu.Lock()
				// ids[(i*number_of_ids)+j] = id
				ids[j] = f.MustNew()
				// mu.Unlock()
			}
			mu.Lock()
			all = append(all, ids...)
			mu.Unlock()

		}()
	}
	wg.Wait()

	for i, id0 := range all {
		for j, id1 := range all {
			if i == j {
				continue
			}
			if id0.String() == id1.String() {
				t.Errorf("ids must be unique: index0: %v index1: %v id:%v", i, j, id0)
				break
			}
		}
	}

}

func TestNew(t *testing.T) {
	var id ID
	var err error
	id, err = New()
	if err != nil {
		t.Fatalf("Default NewFactory default should never return a factory that errors. err %v", err)
	}

	if len(id.String()) != 26 {
		t.Errorf("Default NewFactory should return a valid ulid. id: %v", id)
	}
}

func TestMustNew(t *testing.T) {
	t.Run("dont panic", func(t *testing.T) {
		defer func() {
			if rec := recover(); rec != nil {
				t.Fatalf("MustNew must not panic; panic: %v", rec)
			}
		}()

		id := MustNew()
		if len(id.String()) != 26 {
			t.Errorf("MustNew should return a valid ulid. id: %v", id)
		}

	})

}

func TestMustNewPanic(t *testing.T) {
	t.Run("do panic", func(t *testing.T) {
		f := NewDefaultFactory()
		defer func() {
			if rec := recover(); rec == nil {
				t.Fatal("MustNew must panic if no valid ulid can be generated")
			}
		}()

		var id ID
		id = f.mustNew(&dummyClock{time.Date(1, 1, 1, 12, 30, 1, 1, time.UTC)})
		id = f.mustNew(&dummyClock{time.Date(2263, 1, 1, 12, 30, 1, 1, time.UTC)})

		if id.String() == "" {
			t.Fatalf("id must not be empty")
		}
	})

}
