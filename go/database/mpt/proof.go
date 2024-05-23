package mpt

import (
	"github.com/Fantom-foundation/Carmen/go/common"
)

//go:generate mockgen -source proof.go -destination proof_mocks.go -package witness

// rlpEncodedNode is an RLP encoded MPT node.
type rlpEncodedNode []byte

// proofDb is a database of RLP encoded MPT nodes and their hashes that represent witness proofs.
type proofDb map[common.Hash]rlpEncodedNode

// WitnessProof represents a witness proof.
// It contains a database of MPT nodes and their hashes.
// The proof may be verified by iterating the proof database from the
// root hash to a terminal entry according to a Path
// representing either an account address or a storage key.
// If the proof contains an account or a storage slot for the input address or key,
// the proof is representing this account or storage.
type WitnessProof proofDb

// CreateWitnessProof creates a witness proof for the input account address
// and possibly storage slots of the same account under the input storage keys.
// If the proof cannot be created because either the address or a key does not exist,
// it returns false.
// It may return an error when it occurs in the underlying database.
func CreateWitnessProof(nodeSource NodeSource, root *NodeReference, address common.Address, keys ...common.Key) (WitnessProof, bool, error) {
	return nil, false, nil
}

// Add merges the input witness proof into the current witness proof.
func (p WitnessProof) Add(other WitnessProof) {
	for k, v := range other {
		p[k] = v
	}
}

// ExtractAccount extracts an account and storage from the current witness proof.
// It returns a copy that contains only a proof for the input address and storage keys.
// If the input address or keys are not found in the proof, it returns a proof that
// contains only partial proof consisting only of elements that could be reached
// from the input address and keys in the current proof.
// If the proof cannot be fully extracted, this method returns false.
func (p WitnessProof) ExtractAccount(root common.Hash, address common.Address, keys ...common.Key) (WitnessProof, bool) {
	return WitnessProof{}, false
}

// IsValid checks that the proof has all valid entries.
// It means, RLPs are valid encodings and their hashes match the key in the map.
func (p WitnessProof) IsValid() bool {
	return false
}

// GetAccountInfo extracts an account info from the witness proof for the input root hash and the address.
// If the witness proof contains an account for the input address, it returns its information.
// If the proof does not contain an account, it returns false.
// The method may return an error if the proof is invalid.
func (p WitnessProof) GetAccountInfo(root common.Hash, address common.Address) (AccountInfo, bool, error) {
	path := AddressToNibblePath(address, nil) // TODO hash path without the node source
	node := p.getNext(root, path)
	if node == nil {
		return AccountInfo{}, false, nil
	}
	return node.(*AccountNode).info, true, nil
}

// GetBalance extracts a balance from the witness proof for the input root hash and the address.
// If the witness proof contains an account for the input address, it returns its balance.
// If the proof does not contain an account, it returns false.
// The method may return an error if the proof is invalid.
func (p WitnessProof) GetBalance(root common.Hash, address common.Address) (common.Balance, bool, error) {
	return common.Balance{}, false, nil
}

// GetNonce extracts a nonce from the witness proof for the input root hash and the address.
// If the witness proof contains an account for the input address, it returns its nonce.
// If the proof does not contain an account, it returns false.
// The method may return an error if the proof is invalid.
func (p WitnessProof) GetNonce(root common.Hash, address common.Address) (common.Nonce, bool, error) {
	return common.Nonce{}, false, nil
}

// GetCodeHash extracts a code hash from the witness proof for the input root hash and the address.
// If the witness proof contains an account for the input address, it returns its code hash.
// If the proof does not contain an account, it returns false.
// The method may return an error if the proof is invalid.
func (p WitnessProof) GetCodeHash(root common.Hash, address common.Address) (common.Hash, bool, error) {
	return common.Hash{}, false, nil
}

// GetState extracts a storage slot from the witness proof for the input root hash, account address and the storage key.
// If the witness proof contains an input storage slot for the input key, it returns its value.
// If the proof does not contain a slot, it returns false.
// The method may return an error if the proof is invalid.
func (p WitnessProof) GetState(root common.Hash, address common.Address, key common.Key) (common.Value, bool, error) {
	path := KeyToNibblePath(key, nil) // TODO hash path without the node source
	node := p.getNext(root, path)
	if node == nil {
		return common.Value{}, false, nil
	}
	return node.(*ValueNode).value, true, nil
}

// AllStatesZero checks that all storage slots are empty for the input root hash, account address and the storage key range.
// If the witness proof contains all empty slots for the input key range, it returns true.
// An empty slot is a slot that contains a zero value, or does not exist at all.
// If the proof contains a slot, it returns false.
func (p WitnessProof) AllStatesZero(root common.Hash, address common.Address, from, to common.Key) (bool, error) {
	return false, nil
}

// AllAddressesEmpty checks that all accounts are empty for the input root hash and the address range.
// If the witness proof contains all empty accounts for the input address range, it returns true.
// An empty account is an account that contains a zero balance, nonce, and code hash.
// If the proof contains an account, it returns false.
func (p WitnessProof) AllAddressesEmpty(root common.Hash, from, to common.Address) (bool, error) {
	return false, nil
}

func (p WitnessProof) getNext(root common.Hash, path []Nibble) Node {
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
	//return p.getNext(nextRoot, subPath)

	return nil
}

// proofExtractionVisitor is a visitor that visits MPT nodes and creates a witness proof.
// It hashes and encodes the nodes and stores them into the proof database.
type proofExtractionVisitor struct {
	proof      WitnessProof
	nodeSource NodeSource
	err        error
}

// Visit computes RLP and hash of the visited node and puts it to the proof.
func (p *proofExtractionVisitor) Visit(node Node, _ NodeInfo) VisitResponse {
	data := make([]byte, 0, 1024)
	rlp, err := encodeToRlp(node, p.nodeSource, data)
	if err != nil {
		p.err = err
		return VisitResponseAbort
	}
	hash := common.Keccak256(rlp)

	p.proof[hash] = rlp

	return VisitResponseContinue
}
