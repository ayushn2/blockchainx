package blockchain

import (
	"bytes"

	"github.com/ayushn2/blockchainx.git/wallet"
)

// Since blockchains are public we don't want to save any sensitive information on the blocks

type TxOutput struct{
	Value int
	PubKeyHash []byte
}


// Input just references output
type TxInput struct{
	ID []byte
	Out int
	Signature []byte
	PubKey []byte
}

func NewTxOutput(value int, address string) *TxOutput{
	txo := &TxOutput{value, nil}
	txo.Lock([]byte(address))

	return txo
}

func (in *TxInput) UseskEY(pubKeyHash []byte) bool{
	lockingHash := wallet.PublicKeyHash(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

func (out *TxOutput) Lock(address []byte){
	pubKeyHash :=  wallet.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash) - 4]
	out.PubKeyHash = pubKeyHash
}

func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool{
	return bytes.Compare(out.PubKeyHash,pubKeyHash) == 0 //Checking if the input has the same public key hash as the output
}