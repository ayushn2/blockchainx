package main

import (
	"flag"
	"fmt"
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
	fmt.Println(" add -block BLOCK_DATA - add a block to the chain ")
	fmt.Println(" print - Prints the blocks in the chain")
}

func (cli *CammandLine) validateArgs(){
	if len(os.Args) < 2{
		cli.printUsage()
		runtime.Goexit()//unlike os.Goexit it shutdowns the application by shutting down the go routine
	}
}

func (cli *CammandLine) addBlock(data string){
	cli.blockchain.AddBlock(data)
	fmt.Println("Added Block!")
}

func (cli *CammandLine) printChain(){
	iter := cli.blockchain.Iterator()

	for{
		block := iter.Next()

		fmt.Printf("Prev hash : %x\n",block.PrevHash)
		fmt.Printf("Data : %s\n",block.Data)
		fmt.Printf("Hash : %x\n",block.Hash)
		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n",strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0{
			break
		}
	}
}

func (cli *CammandLine) run(){
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("add",flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print",flag.ExitOnError)
	addBlockData := addBlockCmd.String("block","","Block data")

	switch os.Args[1]{
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	default:
		cli.printUsage()
		runtime.Goexit()
}

	if addBlockCmd.Parsed(){
		if *addBlockData == ""{
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed(){
		cli.printChain()
	}
}

func main(){
	defer os.Exit(0) //To ensure if the go runtime is exited properly
	chain := blockchain.InitBlockChain()

	defer chain.Database.Close()// This only works properly when the go runtime is closed properly

	cli := CammandLine{chain}
	cli.run()
}