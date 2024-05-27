package mpt

import (
	"fmt"
	"github.com/Fantom-foundation/Carmen/go/common"
	"github.com/Fantom-foundation/Carmen/go/database/mpt/rlp"
	"go.uber.org/mock/gomock"
	"slices"
	"testing"
)

func TestDecoder_CanDecodeNodes(t *testing.T) {
	tests := []NodeDesc{
		&Extension{
			path: []Nibble{0x8, 0xe, 0xf},
			next: &Account{},
		},
		&Branch{
			children: Children{
				0x7: &Value{length: 55, key: common.Key{0x1, 0x2, 0x3, 0x4}, value: common.Value{0x4, 0x5, 0x6}},
				0xc: &Value{length: 55, value: common.Value{255}},
			},
		},
		&Value{key: common.Key{0x1, 0x2, 0x3, 0x4}, value: common.Value{0x4, 0x5, 0x6}},
		&Account{
			address: common.Address{0x1, 0x2, 0x3},
			storage: &Value{key: common.Key{0x1, 0x2, 0x3, 0x4}, value: common.Value{0x4, 0x5, 0x6}},
			info:    AccountInfo{Nonce: common.Nonce{0xAA, 0xBB}, Balance: common.Balance{0xCC, 0xDD, 0xEE}, CodeHash: common.Hash{0x11, 0x22, 0x33}},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctxt := newNodeContext(t, ctrl)

			_, shared := ctxt.Build(test)
			handle := shared.GetViewHandle()
			defer handle.Release()
			want := handle.Get()
			rlp, err := encodeToRlp(want, ctxt, []byte{})
			if err != nil {
				t.Fatalf("failed to encode node: %v", err)
			}

			got, err := decodeFromRlp(rlp)
			if err != nil {
				t.Fatalf("failed to decode node: %v", err)
			}

			if got != want {
				t.Errorf("unexpected node, got %v, want %v", got, want)
			}
		})
	}
}

func TestDecoder_CorruptedRlp(t *testing.T) {
	str := rlp.String{Str: []byte{0xFF}}
	threeItemsList := rlp.String{Str: rlp.EncodeInto([]byte{}, rlp.List{Items: []rlp.Item{str, str, str}})}
	longStr := rlp.String{Str: []byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}}
	tooLongNumberItemsList := rlp.String{Str: rlp.EncodeInto([]byte{}, rlp.List{Items: []rlp.Item{longStr, longStr, longStr, longStr}})}

	list := rlp.List{Items: []rlp.Item{str, str, str, str}}

	tests := []struct {
		name string
		rlp  []byte
	}{
		{
			name: "single string",
			rlp:  rlp.EncodeInto([]byte{}, str),
		},
		{
			name: "wrong 3 items list",
			rlp:  rlp.EncodeInto([]byte{}, rlp.List{Items: []rlp.Item{str, str, str}}),
		},
		{ // could be account, but nested inner list has only three items instead of four
			name: "wrong 3 items nested list",
			rlp:  rlp.EncodeInto([]byte{}, rlp.List{Items: []rlp.Item{str, threeItemsList}}),
		},
		{ // could be account, but nested inner list has too long strings
			name: "long strings cannot convert to number",
			rlp:  rlp.EncodeInto([]byte{}, rlp.List{Items: []rlp.Item{str, tooLongNumberItemsList}}),
		},
		{ // could be ext or leaf, but first item is not string
			name: "long strings cannot convert to number",
			rlp:  rlp.EncodeInto([]byte{}, rlp.List{Items: []rlp.Item{list, str}}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := decodeFromRlp(test.rlp); err == nil {
				t.Fatalf("expected error, got nil")
			}
		})
	}
}

func Test_isCompactPathLeafNode(t *testing.T) {
	tests := []struct {
		path   []byte
		isLeaf bool
	}{
		{[]byte{0b_0000_0000}, false},
		{[]byte{0b_0001_0000}, false},
		{[]byte{0b_0010_0000}, true},
		{[]byte{0b_0011_0000}, true},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			if got, want := isCompactPathLeafNode(test.path), test.isLeaf; got != want {
				t.Errorf("unexpected result, got %v, want %v", got, want)
			}
		})
	}
}

func Test_compactPathToNibbles(t *testing.T) {
	tests := []struct {
		path    []byte
		nibbles []Nibble
	}{
		{[]byte{0x00, 0x12, 0x34}, []Nibble{0x1, 0x2, 0x3, 0x4}},      // even
		{[]byte{0x11, 0x23, 0x45}, []Nibble{0x1, 0x2, 0x3, 0x4, 0x5}}, // odd
		{[]byte{0x20, 0x12, 0x34}, []Nibble{0x1, 0x2, 0x3, 0x4}},      // even
		{[]byte{0x31, 0x23, 0x45}, []Nibble{0x1, 0x2, 0x3, 0x4, 0x5}}, // odd
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%x", test), func(t *testing.T) {
			if got, want := compactPathToNibbles(test.path), test.nibbles; !slices.Equal(got, want) {
				t.Errorf("unexpected result, got %v, want %v", got, want)
			}
		})
	}
}
