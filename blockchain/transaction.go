package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

type Transaction struct{
	ID []byte
	Inputs []TxInput
	Outputs []TxOutput
}

// Since blockchains are public we don't want to save any sensitive information on the blocks

type TxOutput struct{
	Value int
	PubKey string
}


// Input just references output
type TxInput struct{
	ID []byte
	Out int
	Sig string
}

func (tx *Transaction) SetID(){
	var encoded bytes.Buffer
	var hash [32]byte

	encode :=gob.NewEncoder(&encoded)
	err := encode.Encode(tx)
	Handle(err)

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func CoinbaseTx(to,data string) *Transaction{
	if data == ""{
		data = fmt.Sprintf("Coins to %s",to)
	}

	txin := TxInput{[]byte{},-1,data}//it has no input and no output
	txout := TxOutput{100,to}

	tx := Transaction{nil,[]TxInput{txin},[]TxOutput{txout}}
	tx.SetID()

	return &tx
}

func NewTransaction(from,to string,amount int,chain *Blockchain) *Transaction{
	var inputs []TxInput
	var outputs []TxOutput

	acc,validOutputs := chain.FindSpendableOutputs(from,amount)

	if acc < amount{
		log.Panic("ERROR : Not enough funds ")
	}

	for txid,outs := range validOutputs{
		txID ,err := hex.DecodeString(txid)
		Handle(err)

		for _,out := range outs{
			input := TxInput{txID,out,from}
			inputs = append(inputs,input)
			//Creating an input for each of the unspent outputs in this transaction
		}
	}

	outputs = append(outputs,TxOutput{amount,to})

	if acc > amount {
		outputs = append(outputs,TxOutput{acc - amount, from})
	}

	tx := Transaction{nil,inputs,outputs}
	tx.SetID()

	return &tx
}

func (tx *Transaction) IsCoinbase() bool{
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

func (in *TxInput) CanUnlock(data string) bool{
	return in.Sig == data
}

func (out *TxOutput) CanBeUnlocked(data string) bool{
	return out.PubKey == data
}