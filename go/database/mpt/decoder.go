package mpt

import (
	"fmt"
	"github.com/Fantom-foundation/Carmen/go/database/mpt/rlp"
	"slices"
)

func decodeFromRlp(data []byte) (Node, error) {
	if slices.Equal(data, emptyStringRlpEncoded) {
		return EmptyNode{}, nil
	}

	item, err := rlp.Decode(data)
	if err != nil {
		return nil, err
	}

	list, ok := item.(*rlp.List)
	if !ok {
		return nil, fmt.Errorf("invalid node type: got: %T, wanted: List", item)
	}

	switch len(list.Items) {
	case 2:
		prefix, ok := list.Items[0].(*rlp.String)
		if !ok {
			return nil, fmt.Errorf("invalid prefix type: got: %T, wanted: String", list.Items[0])
		}
	case 17:
		return decodeBranchNode(list)
	}

	return nil, fmt.Errorf("invalid number of list elements: got: %v, wanted: 2 or 17", len(list.Items))
}

func decodeBranchNode(list *rlp.List) (Node, error) {
	return nil, nil
}
