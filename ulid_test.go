package safeulid

import (
	"fmt"
	"sort"
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

type idss []ID

func (ids idss) Len() int {
	return len(ids)
}

func (ids idss) Less(i, j int) bool {
	return ids[i].ULID.Compare(ids[j].ULID) == -1
}
func (ids idss) Swap(i, j int) {
	ids[i], ids[j] = ids[j], ids[i]
}

func TestNewFactory(t *testing.T) {
	r := constantReader{'x'}
	f := NewFactory(r)
	if f == nil {
		t.Fatalf("NewFactory should never return nil")
	}

	var id ID
	var err error
	id, err = f.new(c)
	if err != nil {
		t.Fatalf("NewFactory default should never return a factory that errors. err %v", err)
	}

	p := 16 // number of parallel coroutines
	number_of_ids := 100
	ids := make(idss, p*number_of_ids)
	for i := 0; i < 16; i++ {
		t.Run(fmt.Sprintf("monotonic-id-%d", i), func(t *testing.T) {
			for j := 0; j < number_of_ids; j++ {
				// 0 * 100 + j 			= 	0
				//							1
				//							...
				//							99
				// 1 * 100 + j 			= 	100
				// 							101
				//							...
				id, err = f.new(c)
				if err != nil {
					t.Errorf("new should not error: %v", err)
				}

				ids[(i*number_of_ids)+j] = id

			}
		})
	}

	sort.Sort(ids)

	for i, id0 := range ids {
		for j, id1 := range ids {
			if i == j {
				continue
			}
			if id0.String() == id1.String() {
				t.Errorf("ids must be unique: index: %v id:%v", i, id0)
				break
			}
		}
	}

	// var ids []ID
	// for i := 0; i < 20; i++ {
	// 	id, err = f.new(c)
	// 	if err != nil {
	// 		t.Errorf("new should not error: %v", err)
	// 	}
	// 	fmt.Printf("%#v\n", id)
	// 	fmt.Printf("%q\n", id.Entropy())
	// 	fmt.Printf("%v\n", id.Time())
	// 	ids = append(ids, id)
	// }

	// t.Errorf("ids: %#v", ids)

	// if len(id.String()) != 26 {
	// 	t.Errorf("Default NewFactory should return a valid ulid. id: %v", id)
	// }

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
		f := NewFactory(nil)
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

// func TestNewFactory(t *testing.T) {
// 	r := constantReader{b: 'x'}
// 	f := NewFactory(r)

// 	if f == nil {
// 		t.Fatalf("NewFactory should never return nil")
// 	}

// }

// type args struct {
// 	entropy io.Reader
// }
// tests := []struct {
// 	name string
// 	args args
// 	want *IDFactory
// }{
// 	// TODO: Add test cases.
// 	{name: "nil entropy", args: args{entropy: nil}, want: N}
// }
// for _, tt := range tests {
// 	t.Run(tt.name, func(t *testing.T) {
// 		if got := NewFactory(tt.args.entropy); !reflect.DeepEqual(got, tt.want) {
// 			t.Errorf("NewFactory() = %v, want %v", got, tt.want)
// 		}
// 	})
// }

// func TestIDFactory_New(t *testing.T) {
// 	type fields struct {
// 		pool sync.Pool
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		want    ID
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			i := &IDFactory{
// 				pool: tt.fields.pool,
// 			}
// 			got, err := i.New()
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("IDFactory.New() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("IDFactory.New() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestIDFactory_MustNew(t *testing.T) {
// 	type fields struct {
// 		pool sync.Pool
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		want   ID
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			i := &IDFactory{
// 				pool: tt.fields.pool,
// 			}
// 			if got := i.MustNew(); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("IDFactory.MustNew() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestMustNew(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		want ID
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := MustNew(); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("MustNew() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestNew(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		want    ID
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := New()
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("New() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_clck_Now(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		c    clck
// 		want time.Time
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := clck{}
// 			if got := c.Now(); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("clck.Now() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_mustNew(t *testing.T) {
// 	type args struct {
// 		t  clock
// 		mr ulid.MonotonicReader
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want ID
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := mustNew(tt.args.t, tt.args.mr); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("mustNew() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_newID(t *testing.T) {
// 	type args struct {
// 		t  clock
// 		mr ulid.MonotonicReader
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    ID
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := newID(tt.args.t, tt.args.mr)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("newID() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("newID() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
