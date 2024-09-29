package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ayushn2/blockchainx.git/wallet"
)

type Transaction struct{
	ID []byte
	Inputs []TxInput
	Outputs []TxOutput
}

func (tx Transaction) Serialize() []byte{
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	Handle(err)

	return encoded.Bytes()
}

func (tx *Transaction) Hash() []byte{
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
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

	txin := TxInput{[]byte{},-1,nil,[]byte(data)}//it has no input and no output
	txout := NewTxOutput(1000,to)// Signing 100 tokens to the account that mines genesis block

	tx := Transaction{nil,[]TxInput{txin},[]TxOutput{*txout}}
	tx.SetID()

	return &tx
}

func NewTransaction(from, to string,amount int,chain *Blockchain) *Transaction{
	
	var inputs []TxInput
	var outputs []TxOutput

	wallets, err := wallet.CreateWallets()
	Handle(err)
	w := wallets.GetWallet(from)
	pubKeyHash := wallet.PublicKeyHash(w.PublicKey)

	acc,validOutputs := chain.FindSpendableOutputs(pubKeyHash,amount)
	
	if acc < amount{
		log.Panic("ERROR : Not enough funds ")
	}
	
	for txid,outs := range validOutputs{
		txID ,err := hex.DecodeString(txid)
		Handle(err)

		for _,out := range outs{
			input := TxInput{txID,out,nil,w.PublicKey}
			inputs = append(inputs,input)
			//Creating an input for each of the unspent outputs in this transaction
		}
	}

	

	outputs = append(outputs,*NewTxOutput(amount,to))

	if acc > amount {
		outputs = append(outputs,*NewTxOutput(acc - amount, from))
	}

	

	tx := Transaction{nil,inputs,outputs}
	
	tx.ID = tx.Hash()
	
	chain.SignTransaction(&tx,w.PrivateKey)
	fmt.Println("Success new transaction")
	return &tx
}

func (tx *Transaction) IsCoinbase() bool{
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction){
	if tx.IsCoinbase(){
		return
	}

	fmt.Println("Private key check:", privKey)

	for _, in := range tx.Inputs{
		
		if prevTXs[hex.EncodeToString(in.ID)].ID == nil{
			log.Panic("ERROR: Previous transaction does not exist")
		}
	}

	txCopy := tx.TrimmedCopy()
	
	for inId, in := range txCopy.Inputs {
		prevTX := prevTXs[hex.EncodeToString(in.ID)]
		txCopy.Inputs[inId].Signature = nil
		txCopy.Inputs[inId].PubKey = prevTX.Outputs[in.Out].PubKeyHash
		fmt.Println("Computing the hash of the transaction")
		txCopy.ID = txCopy.Hash()
		fmt.Printf("Transaction copy ID: %x\n", txCopy.ID)
		txCopy.Inputs[inId].PubKey = nil
		
		fmt.Println("Reached the signing code")
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		fmt.Println("This should be printed after signing attempt")
		if err != nil {
			log.Printf("Error signing transaction: %v", err)
		} else {
			fmt.Println("Successfully signed the transaction")
		}
		signature := append(r.Bytes(), s.Bytes()...)
		
		tx.Inputs[inId].Signature = signature

	}
}

func (tx *Transaction) TrimmedCopy() Transaction{
	var inputs []TxInput
	var outputs []TxOutput

	for _, in := range tx.Inputs{
		inputs = append(inputs, TxInput{in.ID, in.Out, nil, nil})
	}

	for _,out := range tx.Outputs {
		outputs = append(outputs, TxOutput{out.Value, out.PubKeyHash})
	}

	txCopy := Transaction{tx.ID,inputs,outputs}
	return txCopy
}

func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool{
	if tx.IsCoinbase(){
		return true
	}

	for _,in := range tx.Inputs{
		if prevTXs[hex.EncodeToString(in.ID)].ID == nil{
			log.Panic("Previous transaction does not exist")
		}
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inId, in := range tx.Inputs{
		prevTx := prevTXs[hex.EncodeToString(in.ID)]
		txCopy.Inputs[inId].Signature = nil
		txCopy.Inputs[inId].PubKey = prevTx.Outputs[in.Out].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Inputs[inId].PubKey = nil

		r := big.Int{}
		s := big.Int{}
		sigLen := len(in.Signature)
		r.SetBytes(in.Signature[:(sigLen / 2)])
		s.SetBytes(in.Signature[(sigLen/2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(in.PubKey)
		x.SetBytes(in.PubKey[:(keyLen/2)])
		y.SetBytes(in.PubKey[(keyLen/2):])

		rawPubKey := ecdsa.PublicKey{
			Curve: curve,
			X:     &x,
			Y:     &y,
		}
		if ecdsa.Verify(&rawPubKey,txCopy.ID,&r,&s) == false{
			return false
		}

	}

	return true
}

func (tx Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.ID))
	for i, input := range tx.Inputs {
		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:     %x", input.ID))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Out))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PubKey))
	}

	for i, output := range tx.Outputs {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Value))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}