package mpt

import (
	"github.com/Fantom-foundation/Carmen/go/common"
)

//go:generate mockgen -source proof.go -destination proof_mocks.go -package witness

// RlpEncodedNode is an RLP encoded MPT node.
type RlpEncodedNode []byte

// proofDb is a database of RLP encoded MPT nodes and their hashes that represent witness proofs.
type proofDb map[common.Hash]RlpEncodedNode

// WitnessProof represents a witness proof.
// It contains a database of MPT nodes and their hashes.
// The proof may be verified by iterating the proof database from the
// root hash to a terminal entry according to a Path
// representing either an account address or a storage key.
// The proof is valid if it can be iterated to an expected terminal entry.
type WitnessProof proofDb

// Merge merges the input witness proof into the current witness proof.
func (p *WitnessProof) Merge(other *WitnessProof) {
	for k, v := range *other {
		(*p)[k] = v
	}
}

// ExtractAccountProof extracts a sub-proof from the current witness proof.
// It returns a copy that contains only proof for the input address and storage keys.
// If the input address or keys are not found in the current proof, it returns an empty proof.
func (p *WitnessProof) ExtractAccountProof(nodeSource NodeSource, root common.Hash, address common.Address, keys ...common.Key) *WitnessProof {
	return &WitnessProof{}
}

// ExtractStorageProof extracts a sub-proof from the current witness proof.
// It returns a copy that contains only proof for the storage key.
// If the key is not found in the current proof, it returns an empty proof.
func (p *WitnessProof) ExtractStorageProof(nodeSource NodeSource, storageRoot common.Hash, address common.Address, key common.Key) *WitnessProof {
	return &WitnessProof{}
}

// ProveAccount verifies the witness proof against the input root hash and address.
// If the witness proof represents the input account for the input address, it returns the account node.
// If the proof is invalid or does not represent the account, it returns nil.
func (p *WitnessProof) ProveAccount(nodeSource NodeSource, root common.Hash, address common.Address) *AccountNode {
	path := AddressToNibblePath(address, nodeSource)
	node := p.provePath(root, path)
	if node == nil {
		return nil
	}
	return node.(*AccountNode)
}

// ProveStorage verifies the witness proof against the input storage root hash, and storage key.
// If the witness proof represents the input storage slot for the input key, it returns the value node.
// If the proof is invalid or does not represent the storage slot, it returns nil.
func (p *WitnessProof) ProveStorage(nodeSource NodeSource, storageRoot common.Hash, key common.Key) *ValueNode {
	path := KeyToNibblePath(key, nodeSource)
	node := p.provePath(storageRoot, path)
	if node == nil {
		return nil
	}
	return node.(*ValueNode)
}

func (p *WitnessProof) provePath(root common.Hash, path []Nibble) Node {
	// TODO draft of the algorithm

	//rlp, exists := (*p)[root]
	//if !exists {
	//	return nil
	//}
	//node, err := decode(rlp)
	//if err != nil {
	//	return nil // malformed RLP in proof
	//}
	//nextRoot, subPath := getNext(node, path)
	//if nextRoot == empty {
	//	return node
	//}
	//
	//return p.provePath(nextRoot, subPath)

	return nil
}

// CreateWitnessProof creates a witness proof for the input account address
// and possibly storage slots of the same account under the input storage keys.
func CreateWitnessProof(nodeSource NodeSource, root *NodeReference, address common.Address, keys ...common.Key) (*WitnessProof, error) {
	return nil, nil
}

// proofExtractionVisitor is a visitor that visits MPT nodes and creates a witness proof.
// It hashes and encodes the nodes and stores them into the proof database.
type proofExtractionVisitor struct {
	proof      *WitnessProof
	nodeSource NodeSource
	err        error
}

// Visit computes RLP and hash of the visited node and puts it to the proof.
func (p *proofExtractionVisitor) Visit(node Node, _ NodeInfo) VisitResponse {
	data := make([]byte, 0, 1024)
	rlp, err := encode(node, p.nodeSource, data)
	if err != nil {
		p.err = err
		return VisitResponseAbort
	}
	hash := common.Keccak256(rlp)

	(*p.proof)[hash] = rlp

	return VisitResponseContinue
}
