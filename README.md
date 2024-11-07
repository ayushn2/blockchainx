# Go Blockchain with Proof of Work and Command-Line Interface  

This project is a simple blockchain implementation written in Go. It includes functionality for creating blocks, adding transactions, implementing a Proof-of-Work (PoW) consensus mechanism, and interacting with the blockchain via a Command-Line Interface (CLI).  

## Features  

### Blockchain Core Features  

- **Blockchain Structure**:  
  A chain of blocks linked together through cryptographic hashes.  

- **Transactions**:  
  Transactions are stored in each block, showcasing the core idea of a blockchain as a ledger.  

- **Proof of Work (PoW)**:  
  A computational challenge that miners must solve to create a valid block.  

- **Block Details**:  
  Each block contains:  
  - `Index`: The block's position in the chain.  
  - `Timestamp`: When the block was created.  
  - `Transactions`: Data stored in the block.  
  - `Nonce`: A value adjusted to meet the PoW requirements.  
  - `Previous Hash`: The hash of the previous block in the chain.  
  - `Current Hash`: The hash representing the block's contents.  

### CLI Features  

The `cli.go` file provides a command-line interface to interact with the blockchain.  

Commands:  

- **Get Balance**  
    ```bash  
    getbalance -address ADDRESS

Fetch the balance for the given address.

- **Create Blockchain**
    ```bash
    createdblockchain -address ADDRESS

Initializes a new blockchain and rewards the address with the genesis block.

- **Print Blockchain**
    ```bash
    printchain

Prints all blocks in the blockchain.

- **Send Coins**
    ```bash
    send -from FROM -to TO -amount AMOUNT

Creates a transaction and mines a new block.

- **Create Wallet**
    ```bash
    createwallet

Generates a new wallet and address.

- **List Addresses**
    ```bash
    listaddresses

Lists all addresses stored in the wallet.

###Getting Started

Prerequisites

- Go programming language installed (version 1.19 or later recommended).

Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/ayushn2/blockchainx.git  
    cd blockchainx

2. Build the project:
    ```bash
    go build

### Running the CLI

- Start the CLI:
    ```bash
    go run main.go COMMAND [OPTIONS]

Replace COMMAND with one of the CLI commands mentioned above.

### Example Commands

1.	Create a Blockchain:
     ```bash
     go run main.go createblockchain -address your-address

2. Print the Blockchain:
     ```bash
     go run main.go printchain

3. Send Coins:
     ```bash
     go run main.go send -from sender-address -to receiver-address -amount 10

4. Get Balance:
     ```bash
     go run main.go getbalance -address your-address

5. Create a Wallet:
     ```bash
     go run main.go createwallet

6. List Addresses:
     ```bash
     go run main.go listaddresses  

### To Do

-	Implement a peer-to-peer network.
- Enhance transaction validation.
- Add a smart contract-like feature for custom logic.

### Contributing

Contributions are welcome! Feel free to submit issues or pull requests.
