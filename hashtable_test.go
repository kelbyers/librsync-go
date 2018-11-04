package librsync

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignatureHashMap_Get(t *testing.T) {
	strong1 := StrongSignatureHashMap{
		Strong: map[string]int{string(1234): 5678},
	}
	type args struct {
		k uint32
	}
	tests := []struct {
		name  string
		h     *SignatureHashMap
		args  args
		want  StrongSignatureHashMap
		want1 bool
	}{
		{name: "ok", h: &SignatureHashMap{
			Weak: map[uint32]StrongSignatureHashMap{1357: strong1}},
			args: args{k: 1357}, want: strong1, want1: true,
		},
		{name: "not found", h: &SignatureHashMap{
			Weak: map[uint32]StrongSignatureHashMap{1357: strong1}},
			args: args{k: 2468}, want: StrongSignatureHashMap{}, want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.h.Get(tt.args.k)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SignatureHashMap.Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SignatureHashMap.Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSignatureHashMap_Set(t *testing.T) {
	type args struct {
		k uint32
		v StrongSignatureHashMap
	}
	tests := []struct {
		name string
		h    SignatureHashMap
		args args
	}{
		{
			name: "ok",
			h:    NewSignatureHashMap(),
			args: args{k: 1234, v: StrongSignatureHashMap{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.Set(tt.args.k, tt.args.v)
			assert.Equal(t, tt.args.v, tt.h.Weak[tt.args.k])
		})
	}
}

func TestSignatureHashMap_UpdateBlock(t *testing.T) {
	type args struct {
		w uint32
		s []byte
		l int
	}
	tests := []struct {
		name string
		h    SignatureHashMap
		args args
		size int
	}{
		{
			name: "empty",
			h:    NewSignatureHashMap(),
			args: args{w: 1234, s: []byte(`empty`), l: 5678},
			size: 1,
		},
		{
			name: "existing",
			h: SignatureHashMap{Weak: map[uint32]StrongSignatureHashMap{
				1234: StrongSignatureHashMap{
					Strong: map[string]int{block2hash([]byte(`oneitem`)): 9876},
				}}},
			size: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			tt.h.UpdateBlock(tt.args.w, tt.args.s, tt.args.l)
			sh, ok := tt.h.Get(tt.args.w)
			assert.True(ok, "Weak hash added")
			pos, ok := sh.Get(tt.args.s)
			assert.True(ok, "Strong hash added")
			assert.Equal(tt.args.l, pos, "Position added for weak:strong")
			assert.Equal(tt.size, len(tt.h.Weak), "Expected number of items")
		})
	}
}

func TestNewSignatureHashMap(t *testing.T) {
	tests := []struct {
		name string
		want SignatureHashMap
	}{
		{name: "ok", want: SignatureHashMap{
			Weak: map[uint32]StrongSignatureHashMap{},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSignatureHashMap()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSignatureHashMap() = %v, want %v", got, tt.want)
			}
			// make sure the internal map is not a nil map
			got.Weak[1234] = StrongSignatureHashMap{}
		})
	}
}
