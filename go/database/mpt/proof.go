package mpt

import (
	"bytes"
	"errors"
	"github.com/Fantom-foundation/Carmen/go/common"
)

//go:generate mockgen -source proof.go -destination proof_mocks.go -package witness

// RlpEncodedNode is an RLP encoded MPT node.
type RlpEncodedNode []byte

// ProofElement is an entry of a witness proof.
// It contains a hash of an MPT node and the node itself
// encoded as an RLP byte slice.
type ProofElement struct {
	hash common.Hash
	node RlpEncodedNode
}

// Equal returns true if the input proof element is equal to the current proof element.
func (p ProofElement) Equal(other ProofElement) bool {
	return p.hash == other.hash && bytes.Equal(p.node, other.node)
}

// Hash returns the hash of the proof element.
func (p ProofElement) Hash() common.Hash {
	return p.hash
}

// Node returns the RLP encoded node of the proof element.
func (p ProofElement) Node() RlpEncodedNode {
	return p.node
}

// proofDb is a database of MPT nodes and their hashes that represent witness proofs.
type proofDb map[common.Hash]RlpEncodedNode

// proofExtractionVisitor is a visitor that visits MPT nodes and creates a witness proof.
// It hashes and encodes the nodes and stores them into the proof database.
type proofExtractionVisitor struct {
	db         proofDb
	nodeSource NodeSource
	rootHash   *common.Hash
	err        error
}

// Visit computes RLP and hash of the visited node and appends it to the proof.
// It is expected that the nodes are provided via the visitor in the same order
// as they are present in the MPT path from root to the leaf to prove.
func (p *proofExtractionVisitor) Visit(node Node, _ NodeInfo) VisitResponse {
	data := make([]byte, 0, 1024)
	rlp, err := encode(node, p.nodeSource, data)
	if err != nil {
		p.err = err
		return VisitResponseAbort
	}
	hash := common.Keccak256(rlp)
	if p.rootHash != nil {
		p.rootHash = &hash
	}

	p.db[hash] = rlp

	return VisitResponseContinue
}

// WitnessProof represents a witness proof.
// It contains a database of MPT nodes and their hashes,
// and the root hash, which is used as a first lookup to the database.
// The proof may be invalid if iterating the proof database from the
// root hash does not lead to an expected terminal entry.
// In this case the database contains a partial proof and
// the valid flag is set to false.
type WitnessProof struct {
	db    proofDb
	root  common.Hash
	valid bool
}

// Node returns the RLP encoded node and a boolean indicating if the node exists in the proof.
func (w WitnessProof) Node(h common.Hash) (RlpEncodedNode, bool) {
	node, exist := w.db[h]
	return node, exist
}

// Root returns the root hash of the witness proof.
func (w WitnessProof) Root() common.Hash {
	return w.root
}

// IsValid returns true if the witness proof is valid.
func (w WitnessProof) IsValid() bool {
	return w.valid
}

// AccountProof is an account and storage witness proof.
// It contains the witness proof of the account and a list of storage proofs for the storage slots.
// The witness proof is created by iterating the MPT from the root node
// down to an account node following the input address.
// Each node encountered on the path is hashed and stored in a proof.
// Furthermore, the storage proofs are created by iterating the MPT from the account node
// down to the leaf node following the storage key.
// The proof contains one account and possibly multiple storage proofs.
type AccountProof struct {
	accountProof WitnessProof
	storage      map[common.Key]WitnessProof
}

// GetAccountProof returns the account proof.
func (p *AccountProof) GetAccountProof() WitnessProof {
	return p.accountProof
}

// GetStorageProof returns the storage proof for the input storage key.
func (p *AccountProof) GetStorageProof(key common.Key) WitnessProof {
	return p.storage[key]
}

// CreateAccountProof creates a witness proof for the input account
// and possibly storage slots of the same account under the input storage keys.
func CreateAccountProof(nodeSource NodeSource, root *NodeReference, address common.Address, keys ...common.Key) (*AccountProof, error) {
	db := make(proofDb)
	accountVisitor := &proofExtractionVisitor{
		nodeSource: nodeSource,
		db:         db,
	}

	var storageVisitError error

	proof := AccountProof{storage: make(map[common.Key]WitnessProof, len(keys))}
	exists, err := VisitPathToAccount(nodeSource, root, address, MakeVisitor(func(node Node, info NodeInfo) VisitResponse {
		if res := accountVisitor.Visit(node, info); res == VisitResponseAbort {
			return VisitResponseAbort
		}
		// if reached account, prove storage keys and terminate.
		switch account := node.(type) {
		case *AccountNode:
			for _, key := range keys {
				storageVisitor := &proofExtractionVisitor{
					nodeSource: nodeSource,
					db:         db,
				}
				exists, err := VisitPathToStorage(nodeSource, &account.storage, key, storageVisitor)
				if err != nil || storageVisitor.err != nil {
					storageVisitError = errors.Join(storageVisitError, storageVisitor.err, err)
					return VisitResponseAbort
				}
				proof.storage[key] = WitnessProof{db, account.storageHash, exists}
			}

			return VisitResponseAbort
		}

		return VisitResponseContinue
	}))

	if err != nil || storageVisitError != nil || accountVisitor.err != nil {
		return nil, errors.Join(storageVisitError, accountVisitor.err, err)
	}

	var rootHash common.Hash
	if accountVisitor.rootHash != nil {
		rootHash = *accountVisitor.rootHash
	}
	proof.accountProof = WitnessProof{db, rootHash, exists}
	return &proof, nil
}
