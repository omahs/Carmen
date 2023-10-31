package cppstate

/*
import (
	"bytes"
	"errors"
	"github.com/Fantom-foundation/Carmen/go/state"
	"testing"

	"github.com/Fantom-foundation/Carmen/go/common"
)

// directUpdateState is an extended version of the State interface adding support for
// triggering and mocking individual state updates. All its additional members are
// private and not intended to be used outside this package.
type directUpdateState interface {
	state.State

	// createAccount creates a new account with the given address.
	createAccount(address common.Address) error

	// deleteAccount deletes the account with the given address.
	deleteAccount(address common.Address) error

	// setBalance provides balance for the input account address.
	setBalance(address common.Address, balance common.Balance) error

	// setNonce updates nonce of the account for the  input account address.
	setNonce(address common.Address, nonce common.Nonce) error

	// setStorage updates the memory slot for the account address (i.e. the contract) and the memory location key.
	setStorage(address common.Address, key common.Key, value common.Value) error

	// setCode updates code of the contract for the input contract address.
	setCode(address common.Address, code []byte) error
}

var (
	address1 = common.Address{0x01}
	address2 = common.Address{0x02}
	address3 = common.Address{0x03}

	key1 = common.Key{0x01}
	key2 = common.Key{0x02}
	key3 = common.Key{0x03}

	val0 = common.Value{0x00}
	val1 = common.Value{0x01}
	val2 = common.Value{0x02}
	val3 = common.Value{0x03}

	balance1 = common.Balance{0x01}
	balance2 = common.Balance{0x02}
	balance3 = common.Balance{0x03}

	nonce1 = common.Nonce{0x01}
	nonce2 = common.Nonce{0x02}
	nonce3 = common.Nonce{0x03}
)

func TestAccountsAreInitiallyUnknown(t *testing.T) {
	runForEachCppConfig(t, func(t *testing.T, state directUpdateState) {
		account_state, _ := state.Exists(address1)
		if account_state != false {
			t.Errorf("Initial account is not unknown, got %v", account_state)
		}
	})
}

func TestAccountsCanBeCreated(t *testing.T) {
	runForEachCppConfig(t, func(t *testing.T, state directUpdateState) {
		state.createAccount(address1)
		account_state, _ := state.Exists(address1)
		if account_state != true {
			t.Errorf("Created account does not exist, got %v", account_state)
		}
	})
}

func TestAccountsCanBeDeleted(t *testing.T) {
	runForEachCppConfig(t, func(t *testing.T, state directUpdateState) {
		state.createAccount(address1)
		state.deleteAccount(address1)
		account_state, _ := state.Exists(address1)
		if account_state != false {
			t.Errorf("Deleted account is not deleted, got %v", account_state)
		}
	})
}

func TestReadUninitializedBalance(t *testing.T) {
	runForEachCppConfig(t, func(t *testing.T, state directUpdateState) {
		balance, err := state.GetBalance(address1)
		if err != nil {
			t.Fatalf("Error fetching balance: %v", err)
		}
		if (balance != common.Balance{}) {
			t.Errorf("Initial balance is not zero, got %v", balance)
		}
	})
}

func TestWriteAndReadBalance(t *testing.T) {
	runForEachCppConfig(t, func(t *testing.T, state directUpdateState) {
		err := state.setBalance(address1, balance1)
		if err != nil {
			t.Fatalf("Error updating balance: %v", err)
		}
		balance, err := state.GetBalance(address1)
		if err != nil {
			t.Fatalf("Error fetching balance: %v", err)
		}
		if balance != balance1 {
			t.Errorf("Invalid balance read, got %v, wanted %v", balance, balance1)
		}
	})
}

func TestReadUninitializedNonce(t *testing.T) {
	runForEachCppConfig(t, func(t *testing.T, state directUpdateState) {
		nonce, err := state.GetNonce(address1)
		if err != nil {
			t.Fatalf("Error fetching nonce: %v", err)
		}
		if (nonce != common.Nonce{}) {
			t.Errorf("Initial nonce is not zero, got %v", nonce)
		}
	})
}

func TestWriteAndReadNonce(t *testing.T) {
	runForEachCppConfig(t, func(t *testing.T, state directUpdateState) {
		err := state.setNonce(address1, nonce1)
		if err != nil {
			t.Fatalf("Error updating nonce: %v", err)
		}
		nonce, err := state.GetNonce(address1)
		if err != nil {
			t.Fatalf("Error fetching nonce: %v", err)
		}
		if nonce != nonce1 {
			t.Errorf("Invalid nonce read, got %v, wanted %v", nonce, nonce1)
		}
	})
}

func TestReadUninitializedSlot(t *testing.T) {
	runForEachCppConfig(t, func(t *testing.T, state directUpdateState) {
		value, err := state.GetStorage(address1, key1)
		if err != nil {
			t.Fatalf("Error fetching storage slot: %v", err)
		}
		if (value != common.Value{}) {
			t.Errorf("Initial value is not zero, got %v", value)
		}
	})
}

func TestWriteAndReadSlot(t *testing.T) {
	runForEachCppConfig(t, func(t *testing.T, state directUpdateState) {
		err := state.setStorage(address1, key1, val1)
		if err != nil {
			t.Fatalf("Error updating storage: %v", err)
		}
		value, err := state.GetStorage(address1, key1)
		if err != nil {
			t.Fatalf("Error fetching storage slot: %v", err)
		}
		if value != val1 {
			t.Errorf("Invalid value read, got %v, wanted %v", value, val1)
		}
	})
}

func getTestCodeOfLength(size int) []byte {
	res := make([]byte, size)
	for i := 0; i < size; i++ {
		res[i] = byte(i)
	}
	return res
}

func getTestCodes() [][]byte {
	return [][]byte{
		nil,
		{},
		{0xAC},
		{0xAC, 0xDC},
		getTestCodeOfLength(100),
		getTestCodeOfLength(1000),
		getTestCodeOfLength(24577),
	}
}

func TestSetAndGetCode(t *testing.T) {
	runForEachCppConfig(t, func(t *testing.T, state directUpdateState) {
		for _, code := range getTestCodes() {
			err := state.setCode(address1, code)
			if err != nil {
				t.Fatalf("Error setting code: %v", err)
			}
			value, err := state.GetCode(address1)
			if err != nil {
				t.Fatalf("Error fetching code: %v", err)
			}
			if !bytes.Equal(code, value) {
				t.Errorf("Invalid value read, got %v, wanted %v", value, code)
			}
			size, err := state.GetCodeSize(address1)
			if err != nil {
				t.Fatalf("Error fetching code size: %v", err)
			}
			if size != len(code) {
				t.Errorf("Invalid value size read, got %v, wanted %v", size, len(code))
			}
		}
	})
}

func TestSetAndGetCodeHash(t *testing.T) {
	runForEachCppConfig(t, func(t *testing.T, state directUpdateState) {
		for _, code := range getTestCodes() {
			err := state.setCode(address1, code)
			if err != nil {
				t.Fatalf("Error setting code: %v", err)
			}
			hash, err := state.GetCodeHash(address1)
			if err != nil {
				t.Fatalf("Error fetching code: %v", err)
			}
			want := common.GetKeccak256Hash(code)
			if hash != want {
				t.Errorf("Invalid code hash, got %v, wanted %v", hash, want)
			}
		}
	})
}

func initCppStates() []namedStateConfig {
	list := []namedStateConfig{}
	for _, s := range state.GetAllSchemas() {
		list = append(list, []namedStateConfig{
			{"memory", s, castToDirectUpdateState(newCppInMemoryState)},
			{"file", s, castToDirectUpdateState(newCppFileBasedState)},
			{"leveldb", s, castToDirectUpdateState(newCppLevelDbBasedState)},
		}...)
	}
	return list
}

func runForEachCppConfig(t *testing.T, test func(*testing.T, directUpdateState)) {
	for _, config := range initCppStates() {
		config := config
		t.Run(config.name, func(t *testing.T) {
			t.Parallel()
			s, err := config.createState(t.TempDir())
			if err != nil {
				if errors.Is(err, state.UnsupportedConfiguration) {
					t.Skipf("failed to initialize state %s: %v", config.name, err)
				} else {
					t.Fatalf("failed to initialize state %s: %v", config.name, err)
				}
			}
			defer s.Close()
			test(t, s)
		})
	}
}
*/
