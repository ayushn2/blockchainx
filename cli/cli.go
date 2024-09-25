package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/ayushn2/blockchainx.git/blockchain"
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
}

func (cli *CammandLine) validateArgs(){
	if len(os.Args) < 2{
		cli.printUsage()
		runtime.Goexit()//unlike os.Goexit it shutdowns the application by shutting down the go routine
	}
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
		fmt.Println()

		if len(block.PrevHash) == 0{
			break
		}
	}
}

func (cli *CammandLine) createBlockChain(address string){
	// Address will be the person that mines the genesis block
	chain := blockchain.InitBlockChain(address)
	chain.Database.Close()
	fmt.Println("Finished!")
}

func (cli *CammandLine) getBalance(address string){
	chain :=  blockchain.ContinueBlockChain(address)
	defer chain.Database.Close()

	balance := 0
	UTXOs := chain.FindUTXO(address)

	for _,out := range UTXOs{
		balance += out.Value
	}

	fmt.Printf("Balance of %s : %d\n",address,balance)
}

func (cli *CammandLine) send(from,to string, amount int){
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


