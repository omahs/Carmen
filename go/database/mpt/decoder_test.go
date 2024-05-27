package mpt

import (
	"fmt"
	"github.com/Fantom-foundation/Carmen/go/common"
	"github.com/Fantom-foundation/Carmen/go/database/mpt/rlp"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestDecoder_CanDecodeNodes(t *testing.T) {
	tests := []Node{
		&ExtensionNode{},
		&ExtensionNode{path: CreatePathFromNibbles([]Nibble{0x1, 0x2, 0x3}), nextHash: common.Hash{0xA, 0xB, 0xC}},
		&BranchNode{},
		&BranchNode{hashes: [16]common.Hash{EmptyNodeEthereumHash, {0x1, 0x2, 0x3}, EmptyNodeEthereumHash, {0x4, 0x5, 0x6}}},
		&EmptyNode{},
		&ValueNode{key: common.Key{0x1, 0x2, 0x3}, value: common.Value{0x4, 0x5, 0x6}},
		&AccountNode{address: common.Address{0x1, 0x2, 0x3}, storageHash: common.Hash{0x4, 0x5, 0x6},
			info: AccountInfo{Nonce: common.Nonce{0xAA, 0xBB}, Balance: common.Balance{0xCC, 0xDD, 0xEE}, CodeHash: common.Hash{0x11, 0x22, 0x33}}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			nodeSource := NewMockNodeSource(ctrl)
			rlp, err := encodeToRlp(test, nodeSource, []byte{})
			if err != nil {
				t.Fatalf("failed to encode node: %v", err)
			}

			got, err := decodeFromRlp(rlp)
			if err != nil {
				t.Fatalf("failed to decode node: %v", err)
			}

			if got, want := got, test; got != want {
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
