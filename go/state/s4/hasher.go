package s4

//go:generate mockgen -source hasher.go -destination hasher_mocks.go -package s4

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"reflect"

	"github.com/Fantom-foundation/Carmen/go/common"
	"github.com/Fantom-foundation/Carmen/go/state/s4/rlp"
	"golang.org/x/crypto/sha3"
)

// Hasher is an interface for implementations of MPT node hashing algorithms. It is
// intended to be one of the differentiator between S4 and derived schemas.
type Hasher interface {
	// GetHash requests a hash value for the given node. To compute the node's hash,
	// implementations may recursively resolve hashes for other nodes using the given
	// HashSource implementation. Due to its recursive nature, multiple calls to the
	// function may be nested and/or processed concurrently. Thus, implementations are
	// required to be reentrant and thread-safe.
	GetHash(Node, NodeSource, HashSource) (common.Hash, error)
}

type HashSource interface {
	getHashFor(NodeId) (common.Hash, error)
}

// DirectHasher implements a simple, direct node-value hashing algorithm that combines
// the content of individual nodes with the hashes of referenced child nodes into a
// hash for individual nodes.
type DirectHasher struct{}

// GetHash implements the DirectHasher's hashing algorithm.
func (h DirectHasher) GetHash(node Node, _ NodeSource, source HashSource) (common.Hash, error) {
	hash := common.Hash{}
	if _, ok := node.(EmptyNode); ok {
		return hash, nil
	}
	hasher := sha256.New()
	switch node := node.(type) {
	case *AccountNode:
		hasher.Write([]byte{'A'})
		hasher.Write(node.address[:])
		hasher.Write(node.info.Balance[:])
		hasher.Write(node.info.Nonce[:])
		hasher.Write(node.info.CodeHash[:])
		if hash, err := source.getHashFor(node.storage); err == nil {
			hasher.Write(hash[:])
		} else {
			return hash, err
		}

	case *BranchNode:
		hasher.Write([]byte{'B'})
		// TODO: compute sub-tree hashes in parallel
		for _, child := range node.children {
			if hash, err := source.getHashFor(child); err == nil {
				hasher.Write(hash[:])
			} else {
				return hash, err
			}
		}

	case *ExtensionNode:
		hasher.Write([]byte{'E'})
		hasher.Write(node.path.path[:])
		if hash, err := source.getHashFor(node.next); err == nil {
			hasher.Write(hash[:])
		} else {
			return hash, err
		}

	case *ValueNode:
		hasher.Write([]byte{'V'})
		hasher.Write(node.key[:])
		hasher.Write(node.value[:])

	default:
		return hash, fmt.Errorf("unsupported node type: %v", reflect.TypeOf(node))
	}
	hasher.Sum(hash[0:0])
	return hash, nil
}

// Based on Appendix D of https://ethereum.github.io/yellowpaper/paper.pdf
type MptHasher struct{}

// GetHash implements the MPT hashing algorithm.
func (h MptHasher) GetHash(node Node, nodes NodeSource, hashes HashSource) (common.Hash, error) {
	data, err := encode(node, nodes, hashes)
	if err != nil {
		return common.Hash{}, err
	}
	return keccak256(data), nil
}

func keccak256(data []byte) common.Hash {
	return common.GetHash(sha3.NewLegacyKeccak256(), data)
}

func encode(node Node, nodes NodeSource, hashes HashSource) ([]byte, error) {
	switch trg := node.(type) {
	case EmptyNode:
		return encodeEmpty(trg, nodes, hashes)
	case *AccountNode:
		return encodeAccount(trg, nodes, hashes)
	case *BranchNode:
		return encodeBranch(trg, nodes, hashes)
	case *ExtensionNode:
		return encodeExtension(trg, nodes, hashes)
	case *ValueNode:
		return encodeValue(trg, nodes, hashes)
	default:
		return nil, fmt.Errorf("unsupported node type: %v", reflect.TypeOf(node))
	}
}

var emptyStringRlpEncoded = rlp.Encode(rlp.String{})

func encodeEmpty(EmptyNode, NodeSource, HashSource) ([]byte, error) {
	return emptyStringRlpEncoded, nil
}

func encodeBranch(node *BranchNode, nodes NodeSource, hashes HashSource) ([]byte, error) {
	children := node.children
	items := make([]rlp.Item, len(children)+1)

	for i, child := range children {
		if child.IsEmpty() {
			items[i] = rlp.String{}
			continue
		}

		node, err := nodes.getNode(child)
		if err != nil {
			return nil, err
		}

		encoded, err := encode(node, nodes, hashes)
		if err != nil {
			return nil, err
		}

		if len(encoded) >= 32 {
			hash, err := hashes.getHashFor(child)
			if err != nil {
				return nil, err
			}
			encoded = hash[:]
		}
		items[i] = rlp.String{Str: encoded}
	}

	// There is one 17th entry which would be filled if this node is a terminator. However,
	// branch nodes are never terminators in State or Storage Tries.
	items[len(children)] = &rlp.String{}

	var buffer bytes.Buffer
	rlp.List{Items: items}.Write(&buffer)
	return buffer.Bytes(), nil
}

func encodeExtension(node *ExtensionNode, nodes NodeSource, hashes HashSource) ([]byte, error) {
	items := make([]rlp.Item, 2)

	numNibbles := node.path.Length()
	packedNibbles := node.path.GetPackedNibbles()
	items[0] = &rlp.String{Str: encodePartialPath(packedNibbles, numNibbles, false)}

	next, err := nodes.getNode(node.next)
	if err != nil {
		return nil, err
	}
	encoded, err := encode(next, nodes, hashes)
	if err != nil {
		return nil, err
	}
	if len(encoded) >= 32 {
		hash, err := hashes.getHashFor(node.next)
		if err != nil {
			return nil, err
		}
		encoded = hash[:]
	}
	items[1] = &rlp.String{Str: encoded}

	var buffer bytes.Buffer
	rlp.List{Items: items}.Write(&buffer)
	return buffer.Bytes(), nil
}

func encodeAccount(node *AccountNode, nodes NodeSource, hashes HashSource) ([]byte, error) {
	storageRoot := node.storage
	storageHash, err := hashes.getHashFor(storageRoot)
	if err != nil {
		return nil, err
	}

	// Encode the account information to get the value.
	info := node.info
	items := make([]rlp.Item, 4)
	items[0] = &rlp.Uint64{Value: info.Nonce.ToUint64()}
	items[1] = &rlp.BigInt{Value: info.Balance.ToBigInt()}
	items[2] = &rlp.String{Str: storageHash[:]}
	items[3] = &rlp.String{Str: info.CodeHash[:]}
	value := rlp.Encode(rlp.List{Items: items})

	// Encode the leaf node by combining the partial path with the value.
	items = items[0:2]
	items[0] = &rlp.String{Str: encodePath(node.address[:], int(node.pathLength))}
	items[1] = &rlp.String{Str: value}
	return rlp.Encode(rlp.List{Items: items}), nil
}

func encodeValue(node *ValueNode, nodes NodeSource, hashSource HashSource) ([]byte, error) {
	items := make([]rlp.Item, 2)

	// The first item is an encoded path fragment.
	items[0] = &rlp.String{Str: encodePath(node.key[:], int(node.pathLength))}

	// The second item is the value without leading zeros.
	value := node.value[:]
	for len(value) > 0 && value[0] == 0 {
		value = value[1:]
	}
	items[1] = &rlp.String{Str: rlp.Encode(&rlp.String{Str: value[:]})}

	var buffer bytes.Buffer
	rlp.List{Items: items}.Write(&buffer)
	return buffer.Bytes(), nil
}

func encodePath(unhashed []byte, numNibbles int) []byte {
	path := keccak256(unhashed)
	return encodePartialPath(path[32-(numNibbles/2+numNibbles%2):], numNibbles, true)
}

// Requires packedNibbles to include nibbles as [0a bc de] or [ab cd ef]
func encodePartialPath(packedNibbles []byte, numNibbles int, targetsValue bool) []byte {
	// Path encosing derived from Ethereum.
	// see https://github.com/ethereum/go-ethereum/blob/v1.12.0/trie/encoding.go#L37
	oddLength := false
	if numNibbles%2 == 1 {
		oddLength = true
	}

	compact := make([]byte, numNibbles/2+1)

	// The high nibble of the first byte encodes the 'is-value' mark
	// and whether the length is even or odd.
	if targetsValue {
		compact[0] |= 1 << 5
	}
	compact[0] |= (byte(numNibbles) % 2) << 4 // odd flag

	// If there is an odd number of nibbles, the first is included in the
	// low-part of the compact path encoding.
	if oddLength {
		compact[0] |= packedNibbles[0] & 0xf
		packedNibbles = packedNibbles[1:]
	}
	// The rest of the nibbles can be copied.
	copy(compact[1:], packedNibbles)
	return compact
}
