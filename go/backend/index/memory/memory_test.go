package memory

import (
	"github.com/Fantom-foundation/Carmen/go/backend/index"
	"github.com/Fantom-foundation/Carmen/go/common"
	"testing"
)

var (
	A = common.Address{0x01}
	B = common.Address{0x02}
)

func TestImplements(t *testing.T) {
	var memory Memory[*common.Address]
	var _ index.Index[*common.Address, uint64] = &memory
}

func TestStoringIntoMemoryIndex(t *testing.T) {
	memory := NewMemory[common.Address](common.AddressSerializer{})
	defer memory.Close()

	indexA, err := memory.GetOrAdd(A)
	if err != nil {
		t.Errorf("failed add of address A; %s", err)
		return
	}
	if indexA != 0 {
		t.Errorf("first inserted is not 0")
		return
	}
	indexB, err := memory.GetOrAdd(B)
	if err != nil {
		t.Errorf("failed add of address B; %s", err)
		return
	}
	if indexB != 1 {
		t.Errorf("second inserted is not 1")
		return
	}

	if !memory.Contains(A) {
		t.Errorf("memory does not contains inserted A")
		return
	}
	if !memory.Contains(B) {
		t.Errorf("memory does not contains inserted B")
		return
	}

	indexA2, err := memory.GetOrAdd(A)
	if err != nil {
		t.Errorf("failed second add of address A; %s", err)
		return
	}
	if indexA != indexA2 {
		t.Errorf("assigned two different indexes for the same address")
		return
	}

	indexB2, err := memory.GetOrAdd(B)
	if err != nil {
		t.Errorf("failed second add of address B; %s", err)
		return
	}
	if indexB != indexB2 {
		t.Errorf("assigned two different indexes for the same address")
		return
	}
}

func TestHash(t *testing.T) {
	memory := NewMemory[common.Address](common.AddressSerializer{})
	defer memory.Close()

	// the hash is the default one first
	h0 := memory.GetStateHash()

	if (h0 != common.Hash{}) {
		t.Errorf("The hash does not match the default one")
	}

	// the hash must change when adding a new item
	_, _ = memory.GetOrAdd(A)
	h1 := memory.GetStateHash()

	if h0 == h1 {
		t.Errorf("The hash has not changed")
	}

	// the hash remains the same when getting an existing item
	_, _ = memory.GetOrAdd(A)
	h2 := memory.GetStateHash()

	if h1 != h2 {
		t.Errorf("The hash has changed")
	}
}
