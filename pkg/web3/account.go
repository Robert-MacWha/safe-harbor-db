package web3

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"sync/atomic"

	gethCrypto "github.com/ethereum/go-ethereum/crypto"
)

// Account contains a private and public key for a given account.
type Account struct {
	PrivateKey ecdsa.PrivateKey
	Address    Address
	ChainID    int
	nonce      int64
}

type nonceAtClient interface {
	NonceAt(account Address, blockNumber *BigInt) (uint64, error)
}

// NewAccount creates a new account instance from a private key's bytes,
// generating the public key automatically.
// Note: nonce is not set.
func NewAccount(privateKeyBytes []byte, chainID int) (*Account, error) {
	privateKey, err := gethCrypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("error converting private key bytes to ECDSA: %v", err)
	}

	// Derive the public key from the private key
	publicKey := privateKey.PublicKey

	// Derive the Ethereum address from the public key
	commonAddress := gethCrypto.PubkeyToAddress(publicKey)
	address := *CommonToAddress(commonAddress)

	return &Account{
		*privateKey,
		address,
		chainID,
		-1,
	}, nil
}

// NewAccountWithNonce creates a new account instance from a private key's bytes,
// generating the public key automatically and fetching the account's current
// nonce from the uri.
func NewAccountWithNonce(
	privateKeyBytes []byte, chainID int, nodeClient nonceAtClient,
) (*Account, error) {
	account, err := NewAccount(privateKeyBytes, chainID)
	if err != nil {
		return nil, err
	}

	nonce, err := nodeClient.NonceAt(account.Address, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting nonce: %v", err)
	}
	account.nonce = int64(nonce)

	return account, nil
}

// HasAddress returns true if the account has the given address.
func (a *Account) HasAddress(addr Address) bool {
	return bytes.Equal(
		a.Address.ToBytes(), addr.ToBytes(),
	)
}

// SetNonce sets the account's nonce.
func (a *Account) SetNonce(nonce int64) {
	atomic.StoreInt64(&a.nonce, nonce)
}

// GetNonce returns the account's nonce.
func (a *Account) GetNonce() int64 {
	return atomic.LoadInt64(&a.nonce)
}

// IncrementNonce increments the account's nonce by incrementAmount
func (a *Account) IncrementNonce(incrementAmount int64) {
	atomic.AddInt64(&a.nonce, incrementAmount)
}
