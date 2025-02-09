// Copyright (c) 2024 Fantom Foundation
//
// Use of this software is governed by the Business Source License included
// in the LICENSE file and at fantom.foundation/bsl11.
//
// Change Date: 2028-4-16
//
// On the date above, in accordance with the Business Source License, use of
// this software will be governed by the GNU Lesser General Public License v3.

package carmen

import (
	"fmt"
	"math/big"

	"github.com/holiman/uint256"
)

// Amount is a 256-bit unsigned integer used for token values like balances.
type Amount struct {
	internal uint256.Int
}

// NewAmount creates a new U256 Amount from up to 4 uint64 arguments. The
// arguments are given in the Big Endian order. No argument results in a value of zero.
// The constructor panics if more than 4 arguments are given.
func NewAmount(args ...uint64) Amount {
	if len(args) > 4 {
		panic("too many arguments")
	}
	result := Amount{}
	offset := 4 - len(args)
	for i := 0; i < len(args) && i < len(result.internal); i++ {
		result.internal[3-i-offset] = args[i]
	}
	return result
}

// NewAmountFromUint256 creates a new amount from an uint256.
func NewAmountFromUint256(value *uint256.Int) Amount {
	return Amount{internal: *value}
}

// NewAmountFromBytes creates a new Amount instance from up to 32 byte arguments.
// The arguments are given in the Big Endian order. No argument results in a
// value of zero. The constructor panics if more than 32 arguments are given.
func NewAmountFromBytes(bytes ...byte) Amount {
	if len(bytes) > 32 {
		panic("too many arguments")
	}
	result := Amount{}
	result.internal.SetBytes(bytes)
	return result
}

// NewAmountFromBigInt creates a new Amount instance from a big.Int.
func NewAmountFromBigInt(b *big.Int) (Amount, error) {
	if b == nil {
		return NewAmount(), nil
	}
	if b.Sign() < 0 {
		return Amount{}, fmt.Errorf("cannot construct Amount from negative big.Int")
	}
	result := uint256.Int{}
	overflow := result.SetFromBig(b)
	if overflow {
		return Amount{}, fmt.Errorf("big.Int has more than 256 bits")
	}
	return Amount{internal: result}, nil
}

// Uint64 returns the amount as an uint64.
func (a Amount) Uint64() uint64 {
	return a.internal.Uint64()
}

// IsZero returns true if the amount is zero.
func (a Amount) IsZero() bool {
	return a.internal.IsZero()
}

// IsUint64 returns true if the amount is representable as an uint64.
func (a Amount) IsUint64() bool {
	return a.internal.IsUint64()
}

// ToBig returns a bigInt version of the amount.
func (a Amount) ToBig() *big.Int {
	return a.internal.ToBig()
}

// String returns the string representation of the amount.
func (a Amount) String() string {
	return a.internal.String()
}
