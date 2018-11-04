package librsync

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrongSignatureHashMap_Get(t *testing.T) {
	strongSig := []byte(`maphash`)
	strongKey := block2hash(strongSig)

	type args struct {
		k []byte
	}
	tests := []struct {
		name  string
		h     *StrongSignatureHashMap
		args  args
		want  int
		want1 bool
	}{
		{
			name: "ok",
			h:    &StrongSignatureHashMap{Strong: map[string]int{strongKey: 1234}},
			args: args{k: strongSig},
			want: 1234, want1: true,
		},
		{
			name: "not found",
			h:    &StrongSignatureHashMap{Strong: map[string]int{strongKey: 1234}},
			args: args{k: []byte(`wrongway`)},
			want: 0, want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.h.Get(tt.args.k)
			if got != tt.want {
				t.Errorf("StrongSignatureHashMap.Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("StrongSignatureHashMap.Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestStrongSignatureHashMap_Set(t *testing.T) {
	type args struct {
		k string
		l int
	}
	tests := []struct {
		name string
		h    StrongSignatureHashMap
		args args
	}{
		{name: "ok", h: newStrongMap(), args: args{k: "1234", l: 5678}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.Set(tt.args.k, tt.args.l)
			assert.Equal(t, tt.args.l, tt.h.Strong[tt.args.k])
		})
	}
}

func Test_newStrongMap(t *testing.T) {
	tests := []struct {
		name string
		want StrongSignatureHashMap
	}{
		{name: "ok",
			want: StrongSignatureHashMap{Strong: map[string]int{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newStrongMap()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newStrongMap() = %v, want %v", got, tt.want)
			}
			// make sure the internal map is not a nil map
			got.Strong["1234"] = 5678
		})
	}
}
