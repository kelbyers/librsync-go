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
		k []byte
		l int
	}
	tests := []struct {
		name string
		h    StrongSignatureHashMap
		args args
	}{
		{name: "ok", h: newStrongMap(), args: args{k: []byte(`something new`), l: 5678}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.Set(tt.args.k, tt.args.l)
			got, _ := tt.h.Get(tt.args.k)
			assert.Equal(t, tt.args.l, got)
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

func TestStrongSignatureHashMapExtremes(t *testing.T) {
	low32 := make([]byte, 32)
	for i := 0; i < 32; i++ {
		low32[i] = byte(i)
	}
	high32 := make([]byte, 32)
	for i := 0; i < 32; i++ {
		high32[i] = byte(0xff - i)
	}
	tests := []struct {
		name  string
		key   []byte
		value int
	}{
		{name: "low", key: low32, value: 2468},
		{name: "high", key: high32, value: 1357},
		{name: "nihongo", key: []byte(`日本語`), value: 32763},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			tKey := []byte(`1234`)
			tValue := 1234

			m := newStrongMap()
			m.Set(tKey, tValue)
			m.Set(tt.key, tt.value)

			tGot, ok := m.Get(tKey)
			assert.True(ok, "tKey in hashmap")
			assert.Equal(tValue, tGot, "hashmap[tKey] == tValue")

			ttGot, ok := m.Get(tt.key)
			assert.True(ok, "tt.key in hashmap")
			assert.Equal(tt.value, ttGot, "hashmap[tt.key] == tt.value")

			nKey := append([]byte{}, tKey...)
			nKey = append(nKey, tt.key...)
			_, ok = m.Get(nKey)
			assert.False(ok, "key should not be found")

			for k, v := range m.Strong {
				bk := []byte(k)
				assert.Contains([][]byte{tKey, tt.key}, bk, "string key maps back to []byte")
				bv, ok := m.Get(bk)
				assert.True(ok, "Get finds bk")
				assert.Equal(v, bv, "Get finds correct bk")
			}
		})
	}
}
