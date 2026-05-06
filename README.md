# Mingo-Blockchain

<p align="left">
  <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go" />
  <img src="https://img.shields.io/badge/Echo-000000?style=for-the-badge&logo=go&logoColor=white" alt="Echo Framework" />
  <img src="https://img.shields.io/badge/Blockchain-333333?style=for-the-badge&logo=chainlink&logoColor=white" alt="Blockchain" />
  <a href="https://github.com/772005himanshu/Mingo-Blockchain"><img src="https://img.shields.io/badge/GitHub-181717?style=for-the-badge&logo=github&logoColor=white" alt="GitHub" /></a>
  <img src="https://img.shields.io/badge/License-MIT-blue.svg?style=for-the-badge" alt="License: MIT" />
</p>

Mingo-Blockchain is a custom-built, lightweight blockchain implementation written entirely in Go from scratch. Rather than forking an existing project, this repository serves as a comprehensive educational and foundational framework demonstrating how a full blockchain node operates.

From a custom peer-to-peer (P2P) networking layer to a stack-based Virtual Machine capable of executing smart contract instructions, Mingo-Blockchain covers all the core pillars of decentralized ledger technology.

---

## 🏗 System Architecture & File Structure

The project is highly modularized, separating networking, consensus, the virtual machine, and cryptography into their own distinct packages.

```text
.
├── Makefile               # Helper commands for building and testing
├── README.md              # Project documentation
├── go.mod                 # Go module dependencies
├── go.sum                 # Module checksums
├── main.go                # Application entry point and network simulation bootstrap
├── api/
│   └── server.go          # RESTful HTTP API built with labstack/echo for external clients
├── core/
│   ├── Hasher.go          # Core hashing abstractions
│   ├── account_state.go   # Account level state modeling
│   ├── block.go           # Block structure, headers, and cryptographic hashing
│   ├── block_test.go      # Tests for block
│   ├── blockchain.go      # Chain state, block appending, and structural validation
│   ├── blockchain_test.go # Tests for blockchain
│   ├── encoding.go        # Encoding structures (Gob)
│   ├── state.go           # Key-value contract state management
│   ├── state_test.go      # Tests for state
│   ├── storage.go         # Storage interfaces
│   ├── transaction.go     # Transaction schema, collections, and NFT minting structures
│   ├── transaction_test.go# Tests for transactions
│   ├── validator.go       # Rules for validating incoming blocks
│   ├── vm.go              # Custom Stack-based Virtual Machine & instruction set
│   └── vm_test.go         # Tests for Virtual Machine
├── crypto/
│   ├── keypair.go         # Elliptic Curve cryptography, private/public keys, and signatures
│   └── keypair_test.go    # Tests for cryptography
├── network/
│   ├── local_transport.go # In-memory transport used for concurrent local testing
│   ├── local_transport_test.go # Tests for local transport
│   ├── message.go         # Message payloads, parsing, and P2P communication protocols
│   ├── rpc.go             # Remote Procedure Call handling between nodes
│   ├── server.go          # The P2P node server coordinating network activities
│   ├── tcp_transport.go   # Real-world TCP network transport layer
│   ├── transport.go       # General transport interfaces
│   ├── txpool.go          # Mempool for unconfirmed transactions
│   └── txpool_test.go     # Tests for transaction pool
├── types/
│   ├── address.go         # Address derivation and typing
│   ├── hash.go            # 32-byte Hash definitions and type conversions
│   ├── list.go            # Custom list structures
│   └── list_test.go       # Tests for lists
└── util/
    ├── assert.go          # Test assertion utilities
    └── random.go          # Cryptographically secure random generation utilities
```

---

## 🔍 Deep Dive into Core Components

### 1. Network Layer (`network/`)
The networking layer abstracts the concept of a `Transport`. It provides two implementations:
- **TCP Transport**: For real distributed nodes communicating over the internet.
- **Local Transport**: An in-memory, channel-based transport that allows booting up complex networks with multiple nodes inside a single `main.go` process for rapid testing.
The `Server` orchestrates peers, broadcasting new transactions, and syncing blocks via custom RPC messages.

### 2. Core Blockchain (`core/`)
This is the heart of the ledger. 
- **Blocks & Transactions**: Transactions are grouped into blocks. The system includes standard transactions as well as built-in custom transactions like `MintTx` (for NFTs) and `CollectionTx`.
- **Blockchain**: Maintains the longest valid chain, rejecting invalid blocks and tracking block heights.

### 3. Virtual Machine (`core/vm.go`)
Mingo-Blockchain features a bespoke, lightweight **Stack-based Virtual Machine** designed to execute compiled smart contract bytecode.
- **Stack Operations**: Supports `PushInt`, `PushByte`, `Pack`.
- **Arithmetic**: Handles basic math (`InstrAdd`, `InstrSub`, `InstrMul`, `InstrDiv`).
- **State Management**: Provides `InstrStore` and `InstrGet` to allow smart contracts to read and write persistent data to the `contractState`.

### 4. API Layer (`api/`)
To interact with the node from the outside world, a REST API is exposed using the Echo framework. It runs concurrently with the P2P server and intercepts HTTP requests, pushing transactions directly into the node's mempool via Go channels.

---

## 🚀 Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) (1.18+ recommended)

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/772005himanshu/Mingo-Blockchain.git
   cd Mingo-Blockchain
   ```

2. **Download dependencies:**
   ```bash
   go mod tidy
   ```

3. **Run the node (using Makefile):**
   ```bash
   make run
   ```
   
   *Alternatively, you can compile the binary using `make build`, or run the entire test suite using `make test`.*

### What happens when you run `main.go`?
By default, the `main.go` script simulates an entire network locally. It spins up:
- **A Local Node** (with the REST API on port `9000`)
- **Remote Nodes** (communicating via TCP/Local transport)
- **A Late Node** (joins the network after a delay to demonstrate block syncing/catching up)
- An automated **Transaction Sender** that continuously creates and pushes smart contract transactions to the Virtual Machine.

---

## 📡 API Endpoints

Once the node is running, you can communicate with it via the API on `localhost:9000`.

**Fetch a Block by Height**
```bash
curl http://localhost:9000/block/1
```

**Fetch a Block by Hash**
```bash
curl http://localhost:9000/block/<block_hash_hex>
```

**Fetch a Transaction**
```bash
curl http://localhost:9000/tx/<transaction_hash_hex>
```

---

## 📝 License

This project is open-source and available under the MIT License.

