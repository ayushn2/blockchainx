package main

import (
	

	"os"
	"github.com/ayushn2/blockchainx.git/cli"
)



func main(){
	defer os.Exit(0) //To ensure if the go runtime is exited properly
	cli := cli.CammandLine{}
	cli.Run()
}