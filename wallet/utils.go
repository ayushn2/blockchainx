package wallet

import (
	"github.com/ayushn2/blockchainx.git/blockchain"
	"github.com/mr-tron/base58"
)

func Base58Encode(input []byte) []byte{
	encode := base58.Encode(input)

	return []byte(encode)//Required to convert into slice of bytes
}

func Base58Decode(input []byte) []byte{
	decode, err := base58.Decode(string(input[:]))
	blockchain.Handle(err)

	return decode//decode is already a slice of bytes
}

//  Base 58 is a derivative of base64 algo , it uses 6 less characters inside of its alphabet "0 O 1 I : /"