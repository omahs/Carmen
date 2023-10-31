package cppstate

import "github.com/Fantom-foundation/Carmen/go/state"

const (
	CppMemory  state.Variant = "cpp-memory"
	CppFile                  = "cpp-file"
	CppLevelDb               = "cpp-ldb"
)

// NewState is the public interface for creating Carmen state instances. If for the
// given parameters a state can be constructed, the resulting state is returned. If
// construction fails, an error is reported. If the requested configuration is not
// supported, the error is an UnsupportedConfiguration error.
func NewState(params state.Parameters) (state.State, error) {
	switch params.Variant {
	case CppMemory:
		return newCppInMemoryState(params)
	case CppFile:
		return newCppFileBasedState(params)
	case CppLevelDb:
		return newCppLevelDbBasedState(params)
	default:
		return state.NewState(params)
	}
}

func GetAllVariants() []state.Variant {
	return []state.Variant{
		state.GoMemory, state.GoFile, state.GoFileNoCache, state.GoLevelDb, state.GoLevelDbNoCache,
		CppMemory, CppFile, CppLevelDb,
	}
}
