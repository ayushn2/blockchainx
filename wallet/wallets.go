package wallet

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ayushn2/blockchainx.git/blockchain"
)

const walletFile = "./tmp/wallets.data"

type Wallets struct{
	Wallets map[string]*Wallet
}


func CreateWallets() (*Wallets, error){
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFile()

	return &wallets,err
}

func (ws Wallets) GetWallet(address string) Wallet{
	return *ws.Wallets[address]
}

func (ws *Wallets) GetAllAddresses() []string{
	var addresses []string

	for address := range ws.Wallets{
		addresses = append(addresses, address)
	}

	return addresses
}

func (ws *Wallets) AddWallet() string{
	wallet := MakeWallet()
	address := fmt.Sprintf("%s",wallet.Address())

	ws.Wallets[address] = wallet

	return address
}

func (ws *Wallets) LoadFile() error {
    // Check if the wallet file exists
    if _, err := os.Stat(walletFile); os.IsNotExist(err) {
        return err
    }

    // Read the content of the wallet file
    fileContent, err := os.ReadFile(walletFile)
    blockchain.Handle(err)

    // Unmarshal the JSON content into the Wallets struct
    err = json.Unmarshal(fileContent, &ws)
    blockchain.Handle(err)

    return nil
}

// We used json encoding instead of an elliptic curve encoding

func (ws *Wallets) SaveFile() {
	jsonData, err := json.Marshal(ws)
	blockchain.Handle(err)

	err = os.WriteFile(walletFile, jsonData, 0644)
	blockchain.Handle(err)
}