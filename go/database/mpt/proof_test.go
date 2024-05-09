package mpt

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Fantom-foundation/Carmen/go/common"
	"github.com/Fantom-foundation/Carmen/go/database/mpt/shared"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestCreateAccountProof_CanGenerateProof(t *testing.T) {
	ctrl := gomock.NewController(t)

	nodes := []Node{&BranchNode{}, &ExtensionNode{}, &AccountNode{address: common.Address{0xA}}}

	nodeSource := mockNodeSource(ctrl, nodes)
	root := NewNodeReference(EmptyId())

	proof, err := CreateAccountProof(nodeSource, &root, common.Address{0xA})
	if err != nil {
		t.Fatalf("failed to create account proof: %v", err)
	}

	// iterate mocked nodes and check if they are present in the proof database
	// under the correct hash and rlp
	for _, node := range nodes {
		want, err := encode(node, nodeSource, []byte{})
		if err != nil {
			t.Fatalf("failed to encode node: %v", err)
		}
		hash := common.Keccak256(want)
		got, exists := proof.GetAccountProof().Node(hash)
		if !exists {
			t.Errorf("node %x not found in the proof", hash)
		}

		if !bytes.Equal(got, want) {
			t.Errorf("unexpected RLP: got %x, want %x", got, want)
		}
	}

	if got, want := len(proof.GetAccountProof().db), len(nodes); got != want {
		t.Fatalf("unexpected proof length: got %d, want %d", got, want)
	}

	if valid := proof.GetAccountProof().IsValid(); !valid {
		t.Errorf("unexpected proof validity: got %v", valid)
	}

	rootNode, err := encode(nodes[0], nodeSource, []byte{})
	if err != nil {
		t.Fatalf("failed to encode node: %v", err)
	}

	if got, want := proof.GetAccountProof().Root(), common.Keccak256(rootNode); got != want {
		t.Errorf("unexpected root: got %x, want %x", got, want)
	}

}

func TestCreateAccountProof_CannotGenerateProof(t *testing.T) {
	ctrl := gomock.NewController(t)

	nodes := []Node{&EmptyNode{}}

	nodeSource := mockNodeSource(ctrl, nodes)
	root := NewNodeReference(EmptyId())

	proof, err := CreateAccountProof(nodeSource, &root, common.Address{0xA})
	if err != nil {
		t.Fatalf("failed to create account proof: %v", err)
	}

	// the proof database must be empty, the invalid node is not added in the database
	if got, want := len(proof.GetAccountProof().db), 0; got != want {
		t.Fatalf("unexpected proof length: got %d, want %d", got, want)
	}

	if valid := proof.GetAccountProof().IsValid(); valid {
		t.Errorf("unexpected proof validity: got %v", valid)
	}
}

func TestCreateStorageProofs_CanGenerateProofs(t *testing.T) {
	ctrl := gomock.NewController(t)

	nodes := []Node{&BranchNode{}, &ExtensionNode{}, &AccountNode{address: common.Address{0xA}}}

	key1 := common.Key{0x1}
	key2 := common.Key{0x2}

	storageNodes1 := []Node{&BranchNode{}, &ExtensionNode{}, &ValueNode{key: key1}}
	storageNodes2 := []Node{&ExtensionNode{}, &BranchNode{}, &BranchNode{}, &ValueNode{key: key2}}
	storageNodes := append(storageNodes1, storageNodes2...)

	nodeSource := mockNodeSource(ctrl, append(nodes, storageNodes...))
	root := NewNodeReference(EmptyId())

	proof, err := CreateAccountProof(nodeSource, &root, common.Address{0xA}, key1, key2)
	if err != nil {
		t.Fatalf("failed to create account proof: %v", err)
	}

	// iterate mocked nodes and check if they are present in the proof database
	// under the correct hash and rlp
	for _, node := range nodes {
		want, err := encode(node, nodeSource, []byte{})
		if err != nil {
			t.Fatalf("failed to encode node: %v", err)
		}
		hash := common.Keccak256(want)
		got, exists := proof.GetAccountProof().Node(hash)
		if !exists {
			t.Errorf("node %x not found in the proof", hash)
		}

		if !bytes.Equal(got, want) {
			t.Errorf("unexpected RLP: got %x, want %x", got, want)
		}
	}

	if valid := proof.GetAccountProof().IsValid(); !valid {
		t.Errorf("unexpected proof validity: got %v", valid)
	}

	rootNode, err := encode(nodes[0], nodeSource, []byte{})
	if err != nil {
		t.Fatalf("failed to encode node: %v", err)
	}

	if got, want := proof.GetAccountProof().Root(), common.Keccak256(rootNode); got != want {
		t.Errorf("unexpected root: got %x, want %x", got, want)
	}

	keyToNodes := map[common.Key][]Node{
		key1: storageNodes1,
		key2: storageNodes2,
	}

	// continue checking storage proofs are present
	for _, key := range []common.Key{key1, key2} {
		for _, node := range keyToNodes[key] {
			want, err := encode(node, nodeSource, []byte{})
			if err != nil {
				t.Fatalf("failed to encode node: %v", err)
			}
			hash := common.Keccak256(want)
			got, exists := proof.GetStorageProof(key).Node(hash)
			if !exists {
				t.Errorf("node %x not found in the proof", hash)
			}

			if !bytes.Equal(got, want) {
				t.Errorf("unexpected RLP: got %x, want %x", got, want)
			}
		}

		rootNode, err := encode(nodes[len(nodes)-1], nodeSource, []byte{})
		if err != nil {
			t.Fatalf("failed to encode node: %v", err)
		}

		if got, want := proof.GetStorageProof(key).Root(), common.Keccak256(rootNode); got != want {
			t.Errorf("unexpected root: got %x, want %x", got, want)
		}

		if valid := proof.GetStorageProof(key).IsValid(); !valid {
			t.Errorf("unexpected proof validity: got %v", valid)
		}
	}
}

func TestCreateStorageProofs_CannotGenerateProofs(t *testing.T) {
	ctrl := gomock.NewController(t)

	nodes := []Node{&BranchNode{}, &ExtensionNode{}, &AccountNode{address: common.Address{0xA}}}
	storageNodes1 := []Node{&BranchNode{}, &ExtensionNode{}, &EmptyNode{}}                // does not lead to a value node
	storageNodes2 := []Node{&ExtensionNode{}, &BranchNode{}, &BranchNode{}, &EmptyNode{}} // does not lead to a value node

	nodeSource := mockNodeSource(ctrl, append(nodes, append(storageNodes1, storageNodes2...)...))
	root := NewNodeReference(EmptyId())

	key1 := common.Key{0x1}
	key2 := common.Key{0x2}
	proof, err := CreateAccountProof(nodeSource, &root, common.Address{0xA}, key1, key2)
	if err != nil {
		t.Fatalf("failed to create proof: %v", err)
	}

	// account proof is present
	// iterate mocked nodes and check if they are present in the proof database
	// under the correct hash and rlp
	for _, node := range nodes {
		want, err := encode(node, nodeSource, []byte{})
		if err != nil {
			t.Fatalf("failed to encode node: %v", err)
		}
		hash := common.Keccak256(want)
		got, exists := proof.GetAccountProof().Node(hash)
		if !exists {
			t.Errorf("node %x not found in the proof", hash)
		}

		if !bytes.Equal(got, want) {
			t.Errorf("unexpected RLP: got %x, want %x", got, want)
		}
	}

	if valid := proof.GetAccountProof().IsValid(); !valid {
		t.Errorf("unexpected proof validity: got %v", valid)
	}

	rootNode, err := encode(nodes[0], nodeSource, []byte{})
	if err != nil {
		t.Fatalf("failed to encode node: %v", err)
	}

	if got, want := proof.GetAccountProof().Root(), common.Keccak256(rootNode); got != want {
		t.Errorf("unexpected root: got %x, want %x", got, want)
	}

	keyToNodes := map[common.Key][]Node{
		key1: storageNodes1,
		key2: storageNodes2,
	}

	// continue checking storage proofs are present
	for _, key := range []common.Key{key1, key2} {
		for i, node := range keyToNodes[key] {
			want, err := encode(node, nodeSource, []byte{})
			if err != nil {
				t.Fatalf("failed to encode node: %v", err)
			}
			hash := common.Keccak256(want)
			got, exists := proof.GetStorageProof(key).Node(hash)
			if !exists {
				if i == len(keyToNodes[key])-1 {
					// last node is not present in the proof as it was invalid
					continue
				}
				t.Errorf("node %x not found in the proof", hash)
			}

			if !bytes.Equal(got, want) {
				t.Errorf("unexpected RLP: got %x, want %x", got, want)
			}
		}

		rootNode, err := encode(nodes[len(nodes)-1], nodeSource, []byte{})
		if err != nil {
			t.Fatalf("failed to encode node: %v", err)
		}

		if got, want := proof.GetStorageProof(key).Root(), common.Keccak256(rootNode); got != want {
			t.Errorf("unexpected root: got %x, want %x", got, want)
		}

		if valid := proof.GetStorageProof(key).IsValid(); valid {
			t.Errorf("unexpected proof validity: got %v", valid)
		}
	}
}

func TestCreateStorageProofs_AccountDoesNotExist(t *testing.T) {
	ctrl := gomock.NewController(t)

	nodes := []Node{&EmptyNode{}}

	nodeSource := mockNodeSource(ctrl, nodes)
	root := NewNodeReference(EmptyId())

	key1 := common.Key{0x1}
	key2 := common.Key{0x2}
	proof, err := CreateAccountProof(nodeSource, &root, common.Address{}, key1, key2)
	if err != nil {
		t.Fatalf("failed to create proof: %v", err)
	}

	// the proof database must be empty, the invalid node is not added in the database
	if got, want := len(proof.GetAccountProof().db), 0; got != want {
		t.Fatalf("unexpected proof length: got %d, want %d", got, want)
	}

	if valid := proof.GetAccountProof().IsValid(); valid {
		t.Errorf("unexpected proof validity: got %v", valid)
	}

	for _, key := range []common.Key{key1, key2} {
		if valid := proof.GetStorageProof(key).IsValid(); valid {
			t.Errorf("unexpected proof validity: got %v", valid)
		}
	}
}

func TestCreateStorageProofs_CannotGenerateProofs_FailingNodeSources(t *testing.T) {
	ctrl := gomock.NewController(t)

	injectedErr := fmt.Errorf("injected error")
	var node Node

	tests := []struct {
		name string
		mock func(*MockNodeSource)
	}{
		{
			name: "call in account proof fails",
			mock: func(mock *MockNodeSource) {
				mock.EXPECT().getViewAccess(gomock.Any()).Return(shared.MakeShared(node).GetViewHandle(), injectedErr)
			},
		},
		{
			name: "call in storage proof fails",
			mock: func(mock *MockNodeSource) {
				var account Node = &AccountNode{address: common.Address{0xA}}
				gomock.InOrder(
					mock.EXPECT().getViewAccess(gomock.Any()).Return(shared.MakeShared(account).GetViewHandle(), nil),
					mock.EXPECT().getViewAccess(gomock.Any()).Return(shared.MakeShared(node).GetViewHandle(), injectedErr),
				)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nodeSource := NewMockNodeSource(ctrl)
			nodeSource.EXPECT().getConfig().AnyTimes().Return(S5LiveConfig)
			mockPathHashing(nodeSource)

			test.mock(nodeSource)
			root := NewNodeReference(EmptyId())

			if _, err := CreateAccountProof(nodeSource, &root, common.Address{0xA}, common.Key{0x1}); !errors.Is(err, injectedErr) {
				t.Errorf("getting proof should fail")
			}
		})
	}
}

func TestCreateAccountProofs_FailingNodeSources(t *testing.T) {
	ctrl := gomock.NewController(t)

	injectedErr := fmt.Errorf("injected error")

	nodeSource := NewMockNodeSource(ctrl)
	nodeSource.EXPECT().getConfig().AnyTimes().Return(S5LiveConfig)
	mockPathHashing(nodeSource)

	var account Node = &AccountNode{address: common.Address{0xA}}
	nodeSource.EXPECT().getViewAccess(gomock.Any()).Return(shared.MakeShared(account).GetViewHandle(), injectedErr)

	root := NewNodeReference(EmptyId())

	if _, err := CreateAccountProof(nodeSource, &root, common.Address{0xA}); !errors.Is(err, injectedErr) {
		t.Errorf("getting proof should fail")
	}
}

func TestCreateStorageProofs_HasherFails(t *testing.T) {
	ctrl := gomock.NewController(t)

	injectedErr := fmt.Errorf("injected error")

	mockHasher := NewMockhasher(ctrl)
	mockHasher.EXPECT().encode(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte{}, injectedErr)
	mockHasherFunc := func() hasher {
		return mockHasher
	}

	nodeSource := NewMockNodeSource(ctrl)
	nodeSource.EXPECT().getConfig().AnyTimes().Return(MptConfig{Hashing: hashAlgorithm{createHasher: mockHasherFunc}})
	mockPathHashing(nodeSource)

	var account Node = &AccountNode{address: common.Address{0xA}}
	nodeSource.EXPECT().getViewAccess(gomock.Any()).Return(shared.MakeShared(account).GetViewHandle(), nil)

	root := NewNodeReference(EmptyId())

	if _, err := CreateAccountProof(nodeSource, &root, common.Address{0xA}, common.Key{0x1}); !errors.Is(err, injectedErr) {
		t.Errorf("getting proof should fail")
	}
}

func mockPathHashing(nodeSource *MockNodeSource) {
	nodeSource.EXPECT().hashKey(gomock.Any()).AnyTimes().DoAndReturn(func(key common.Key) common.Hash {
		var h common.Hash
		copy(h[:], key[:])
		return h
	})
	nodeSource.EXPECT().hashAddress(gomock.Any()).AnyTimes().DoAndReturn(func(address common.Address) common.Hash {
		var h common.Hash
		copy(h[:], address[:])
		return h
	})
}

func mockNodeSource(ctrl *gomock.Controller, nodes []Node) NodeSource {
	nodeSource := NewMockNodeSource(ctrl)
	nodeSource.EXPECT().getConfig().AnyTimes().Return(S5LiveConfig)
	mockPathHashing(nodeSource)

	calls := make([]*gomock.Call, 0, len(nodes))
	for _, node := range nodes {
		call := nodeSource.EXPECT().getViewAccess(gomock.Any()).Return(shared.MakeShared(node).GetViewHandle(), nil)
		calls = append(calls, call)
	}

	gomock.InOrder(calls...)

	return nodeSource
}

//
//func TestCreateProofs_From_Database(t *testing.T) {
//	keys := getTestKeys(50)
//	addresses := getTestAddresses(10)
//	for _, variant := range fileAndMemVariants {
//		for _, config := range []MptConfig{S5LiveConfig, S5ArchiveConfig} {
//			for forestConfigName, forestConfig := range forestConfigs {
//				t.Run(fmt.Sprintf("%s-%s-%s", variant.name, config.Name, forestConfigName), func(t *testing.T) {
//					forest, err := variant.factory(t.TempDir(), config, forestConfig)
//					if err != nil {
//						t.Fatalf("failed to open forest: %v", err)
//					}
//					defer func() {
//						if err := forest.Close(); err != nil {
//							t.Fatalf("cannot close db: %v", err)
//						}
//					}()
//
//					rootRef := NewNodeReference(EmptyId())
//
//					// fill-in forest with data
//					for i, address := range addresses {
//						root, err := forest.SetAccountInfo(&rootRef, address, AccountInfo{Balance: common.Balance{byte(i + 1)}, Nonce: common.Nonce{1}})
//						if err != nil {
//							t.Fatalf("cannot create an account: %v", err)
//						}
//						rootRef = root
//
//						for j, key := range keys {
//							root, err := forest.SetValue(&rootRef, address, key, common.Value{byte(i + 1), byte(j + 1)})
//							if err != nil {
//								t.Fatalf("cannot create an account: %v", err)
//							}
//							rootRef = root
//						}
//
//						// trigger update of dirty hashes
//						// witness proof cannot be computed on a dirty trie
//						if _, _, err := forest.updateHashesFor(&rootRef); err != nil {
//							t.Errorf("failed to compute hash: %v", err)
//						}
//					}
//
//					// pick just one address to prove
//					proveAddress := addresses[0]
//
//					hasher := config.Hashing.createHasher()
//
//					type slotId struct {
//						address common.Address
//						key     common.Key
//					}
//
//					accountNodes := make(map[common.Address][]rlpAndHash)
//					storageNodes := make(map[slotId][]rlpAndHash)
//
//					var stack []rlpAndHash
//					var accountDepth int
//					var currentAccount common.Address
//
//					// Collect nodes from the forest.
//					// It iterates over the trie and collects the nodes
//					// using a stack to keep track of the nodes on the path.
//					// When a leaf is reached, the nodes are collected and stored.
//					// It serves as a reference to compare the proof later.
//					// This way of collecting paths exploits inner implementation details
//					// of the trie visitor that uses depth-first search.
//					if err := forest.VisitTrie(&rootRef, MakeVisitor(func(node Node, info NodeInfo) VisitResponse {
//						rlp, err := encode(node, forest, []byte{})
//						if err != nil {
//							t.Fatalf("failed to encode node: %v", err)
//						}
//						hash := common.Keccak256(rlp)
//						n := rlpAndHash{hash, rlp}
//
//						stack = stack[0:*info.Depth]
//						stack = append(stack, n)
//
//						// collect the discovered path if a leaf node found
//						switch n := node.(type) {
//						case *AccountNode:
//							nodes := make([]rlpAndHash, len(stack))
//							copy(nodes, stack)
//							accountNodes[n.address] = nodes
//							accountDepth = len(stack)
//							currentAccount = n.address
//						case *ValueNode:
//							nodes := make([]rlpAndHash, len(stack)-accountDepth)
//							copy(nodes, stack[accountDepth:])
//							storageNodes[slotId{currentAccount, n.key}] = nodes
//						}
//
//						return VisitResponseContinue
//					})); err != nil {
//						t.Fatalf("failed to visit trie: %v", err)
//					}
//
//					// proof only part of the inserted keys
//					rand.Shuffle(len(keys), func(i, j int) {
//						keys[i], keys[j] = keys[j], keys[i]
//					})
//					proveKeys := keys[0 : len(keys)/2]
//
//					proof, err := CreateAccountProof(forest, &rootRef, proveAddress, proveKeys...)
//					if err != nil {
//						t.Fatalf("failed to create proof: %v", err)
//					}
//
//					convert := func(items []ProofElement) []rlpAndHash {
//						nodes := make([]rlpAndHash, len(items))
//						for i, item := range items {
//							nodes[i] = rlpAndHash{item.Hash(), item.Node()}
//						}
//						return nodes
//					}
//
//					// check the proof matches the collected nodes for an account
//					if got, want := convert(proof.GetAccountProof()), accountNodes[proveAddress]; !slices.EqualFunc(got, want, rlpAndHash.Equal) {
//						t.Errorf("unexpected proof: got %v, want %v", got, want)
//					}
//
//					// check the proof matches the collected nodes for storage
//					for _, key := range proveKeys {
//						if got, want := convert(proof.GetStorageProof(key)), storageNodes[slotId{proveAddress, key}]; !slices.EqualFunc(got, want, rlpAndHash.Equal) {
//							t.Errorf("unexpected proof: got %v, want %v", got, want)
//						}
//					}
//				})
//			}
//		}
//	}
//}

type rlpAndHash struct {
	hash common.Hash
	rlp  RlpEncodedNode
}

func (r rlpAndHash) Equal(other rlpAndHash) bool {
	return r.hash == other.hash && bytes.Equal(r.rlp, other.rlp)
}
