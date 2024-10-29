package web3

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
)

// BigInt acts as a type alias for big.Int that can easily be marshalled and
// unmarshalled as a hexstring
type BigInt struct {
	*big.Int
}

// NewBigInt creates a new web3.BigInt from a *big.Int
func NewBigInt(b *big.Int) *BigInt {
	return &BigInt{b}
}

// NewBigIntFromFloat creates a new web3.BigInt from a float64
func NewBigIntFromFloat(f float64) *BigInt {
	i, _ := big.NewFloat(f).Int(nil)
	return NewBigInt(i)
}

// IntToBig creates a new web3.BigInt from a int64
func IntToBig(b int64) *BigInt {
	return &BigInt{big.NewInt(b)}
}

// Zero just returns a zero BigInt
func Zero() *BigInt {
	return IntToBig(0)
}

// BytesToBig creates a new web3.BigInt from a []byte
func BytesToBig(b []byte) *BigInt {
	i := big.NewInt(0)
	i.SetBytes(b)
	return NewBigInt(i)
}

// ToBigDotInt converts a web3.BigInt to a big.Int
func (b *BigInt) ToBigDotInt() *big.Int {
	var bigInt big.Int
	return bigInt.Set(b.Int)
}

// MarshalJSON override
func (b BigInt) MarshalJSON() ([]byte, error) {
	if b.Int == nil {
		return json.Marshal("0x0")
	}
	hexValue := "0x" + b.Text(16)
	return json.Marshal(hexValue)
}

// Sub sets z to the difference x-y and returns z.
func (b *BigInt) Sub(x, y *BigInt) *BigInt {
	b.Int.Sub(x.Int, y.Int)
	return b
}

// Add sets z to the sum x+y and returns z.
func (b *BigInt) Add(x, y *BigInt) *BigInt {
	b.Int.Add(x.Int, y.Int)
	return b
}

// DirAdd is Direct Addition without pointers
func (b *BigInt) DirAdd(x, y BigInt) BigInt {
	b.Int.Add(x.Int, y.Int)
	return *b
}

// DirSub is Direct Subtraction without pointers
func (b *BigInt) DirSub(x, y BigInt) BigInt {
	b.Int.Sub(x.Int, y.Int)
	return *b
}

// Abs creates a new big int of value |b|
func (b *BigInt) Abs() *BigInt {
	return NewBigInt(new(big.Int).Abs(b.Int))
}

// Mul sets z to the product x*y and returns z.
func (b *BigInt) Mul(x, y *BigInt) *BigInt {
	b.Int.Mul(x.Int, y.Int)
	return b
}

// Div sets z to the quotient x/y and returns z.
func (b *BigInt) Div(x, y *BigInt) *BigInt {
	b.Int.Div(x.Int, y.Int)
	return b
}

// Exp sets z to x**y and returns z.
func (b *BigInt) Exp(x, y *BigInt) *BigInt {
	b.Int.Exp(x.Int, y.Int, nil)
	return b
}

// Cmp compares x and y and returns:
//
//	-1 if x <  y
//	 0 if x == y
//	+1 if x >  y
func (b *BigInt) Cmp(y *BigInt) int {
	return b.Int.Cmp(y.Int)
}

// SubtractIfNotNil subtracts y from x if x and y are not nil
func (b *BigInt) SubtractIfNotNil(x, y *BigInt) *BigInt {
	if x != nil && y != nil {
		b.Int.Sub(x.Int, y.Int)
	}
	return b
}

// UnmarshalJSON override
func (b *BigInt) UnmarshalJSON(data []byte) error {
	var hexValue string
	err := json.Unmarshal(data, &hexValue)
	if err != nil {
		return fmt.Errorf("error unmarshalling BigInt: %v", err)
	}

	if b.Int == nil {
		b.Int = new(big.Int)
	}

	value := strings.TrimPrefix(hexValue, "0x")
	if value == "" {
		return nil
	}

	_, ok := b.SetString(value, 16)
	if !ok {
		return fmt.Errorf("error setting big.Int value to %v", value)
	}

	return nil
}

// UnmarshalYAML override
func (b *BigInt) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return fmt.Errorf("error unmarshalling BigInt: %v", err)
	}

	if b.Int == nil {
		b.Int = new(big.Int)
	}

	_, ok := b.SetString(value, 0)
	if !ok {
		return fmt.Errorf("error setting big.Int value to %v", value)
	}
	return nil
}

// CeilBigFloat returns the result rounded up to the nearest integer
func CeilBigFloat(f *big.Float) *BigInt {
	half := 0.5
	delta := half
	if f.Sign() < 0 {
		delta = -half
	}
	f.Add(f, new(big.Float).SetFloat64(delta))
	bint, _ := f.Int(nil)

	return NewBigInt(bint)
}

// RoundBigFloat rounds a big.Float to the nearest integer
func RoundBigFloat(f *big.Float) *BigInt {
	// Creating a half value as a big.Float
	halfFloat := 0.5
	half := big.NewFloat(halfFloat)

	// If the sign is negative, we subtract 0.5; otherwise, we add 0.5
	if f.Sign() < 0 {
		f.Sub(f, half)
	} else {
		f.Add(f, half)
	}

	// Extract the nearest integer
	bint, _ := f.Int(nil)

	// Convert the result to BigInt (using your NewBigInt function)
	return NewBigInt(bint)
}

// MultiplyWeb3BigIntByFloat multiplies a BigInt by a float
func (b *BigInt) MultiplyWeb3BigIntByFloat(multiplier float64) *BigInt {
	bigInt := b.Int
	bigFloat := new(big.Float).SetInt(bigInt)
	bigFloat = new(big.Float).Mul(bigFloat, big.NewFloat(multiplier))
	res := CeilBigFloat(bigFloat)
	return res
}

// GobEncode helps Gob encode it into bytes
func (b BigInt) GobEncode() ([]byte, error) {
	if b.Int == nil {
		// Encode a nil *big.Int as an empty slice
		return []byte{}, nil
	}
	return b.Int.Bytes(), nil
}

// GobDecode helps Gob decodei it from bytes to bigInt
func (b *BigInt) GobDecode(data []byte) error {
	if b.Int == nil {
		b.Int = new(big.Int)
	}
	b.Int.SetBytes(data)
	return nil
}
