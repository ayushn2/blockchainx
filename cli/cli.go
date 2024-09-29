package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/ayushn2/blockchainx.git/blockchain"
	"github.com/ayushn2/blockchainx.git/wallet"
)

type CammandLine struct{
	blockchain *blockchain.Blockchain
}

func (cli *CammandLine) printUsage(){
	fmt.Println("Usage : ")
	fmt.Println(" getbalance -address ADDRESS - get the balance for the address")
	fmt.Println(" createblockchain -address ADDRESS creates a blockchain")
	fmt.Println(" printchain - Prints the blocks in the chain")
	fmt.Println(" send -from FROM -to TO -amount AMOUNT - Send amount")
	fmt.Println(" createwallet - Creates a new wallet")
	fmt.Println(" listaddresses - Lists the addresses in out wallet file")
}

func (cli *CammandLine) validateArgs(){
	if len(os.Args) < 2{
		cli.printUsage()
		runtime.Goexit()//unlike os.Goexit it shutdowns the application by shutting down the go routine
	}
}

func (cli *CammandLine) listAddresses(){
	wallets, _ := wallet.CreateWallets()
	addresses := wallets.GetAllAddresses()

	for _, address := range addresses{
		fmt.Println(address)
	}
}

func (cli *CammandLine) createWallet(){
	wallets, _ := wallet.CreateWallets()
	address := wallets.AddWallet()
	wallets.SaveFile()

	fmt.Printf("New address is: %s\n",address)
}

func (cli *CammandLine) printChain(){
	chain := blockchain.ContinueBlockChain("")
	defer chain.Database.Close()
	iter := chain.Iterator()

	for{
		block := iter.Next()

		fmt.Printf("Prev hash : %x\n",block.PrevHash)
		fmt.Printf("Hash : %x\n",block.Hash)
		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n",strconv.FormatBool(pow.Validate()))
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}

		fmt.Println()

		if len(block.PrevHash) == 0{
			break
		}
	}
}

func (cli *CammandLine) createBlockChain(address string){

	if !wallet.ValidateAddress(address){
		log.Panic("Address is not valid")
	}

	// Address will be the person that mines the genesis block
	chain := blockchain.InitBlockChain(address)
	chain.Database.Close()
	fmt.Println("Finished!")
}

func (cli *CammandLine) getBalance(address string){

	if !wallet.ValidateAddress(address){
		log.Panic("Address is not valid")
	}

	chain :=  blockchain.ContinueBlockChain(address)
	defer chain.Database.Close()

	balance := 0
	pubKeyHash := wallet.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[ 1: len(pubKeyHash) - 4]


	UTXOs := chain.FindUTXO(pubKeyHash)

	for _,out := range UTXOs{
		balance += out.Value
	}

	fmt.Printf("Balance of %s : %d\n",address,balance)
}

func (cli *CammandLine) send(from,to string, amount int){

	if !wallet.ValidateAddress(from){
		log.Panic("Address is not valid")
	}

	if !wallet.ValidateAddress(to){
		log.Panic("Address is not valid")
	}

	chain := blockchain.ContinueBlockChain(from)
	defer chain.Database.Close()
	tx := blockchain.NewTransaction(from,to,amount,chain)
	
	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Println("Success!") 
}

func (cli *CammandLine) Run(){
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance",flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain",flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send",flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain",flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet",flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listaddresses",flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address","","The address has balance : ")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1]{
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil{
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil{
			log.Panic(err)
		}
	case "listaddresses":
		err := listAddressesCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[:2])
		blockchain.Handle(err)
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil{
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		runtime.Goexit()
}

	if getBalanceCmd.Parsed(){
		if *getBalanceAddress == ""{
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed(){
		if *createBlockchainAddress == ""{
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockChain(*createBlockchainAddress)
	}
	if createWalletCmd.Parsed(){
		cli.createWallet()
	}
	if listAddressesCmd.Parsed(){
		cli.listAddresses()
	}
	if printChainCmd.Parsed() {
		cli.printChain()
	}
	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}


