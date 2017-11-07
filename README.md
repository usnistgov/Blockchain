Overview of the NIST Blockchain Prototypes
------------------------------------------
Stephen Nightingale, NIST Information Technology Laboratory

Bitcoin is the first and the iconic instance of a Blockchain application for a
cryptocurrency. It has the original (nakamoto) implementation in C++, and a
later implementation in [Go][go], available in the [btcsuite][btcd] on GitHub.
In full, it is a cryptocurrency with decentralized control of the ledger, and
universal access. Anyone can participate as a wallet user, getting and spending
Bitcoins, or by deploying a node, tracking and mining Bitcoins and contributing
to the ledger maintenance.

Bitcoin introduced a new data structure, the blockchain, as a secure public
means of storing and transferring currency, and other information. It employs
cryptographic concepts and tools such as public/private key pairs, and hashes
(and hashes of hashes). Since the implementation employs complex layers of
dereferencing, it is difficult code to read, understand and modify. The
object of this project, then, is to develop cryptocurrency models from the
simplest, to the progressively more complex, while introducing concepts using
in the development of Bitcoin.

We hope that this will offer the curious a path to understanding Bitcoin, and
code to use to develop alternative concepts. It is offered as a tool for
modelling Blockchain development and exploration.

All the early models in this repository are centralized. A summary of the
models is:

- Model 0: introduces immutable coins, identified by transaction ID,
  with ownership identified by hash. The ledger is a linear list. The
  outputs are signed by the sender to indicate transfer of ownership.
  [Model0/README.md](Model0/README.md) explains how the original model works.
- Model 1: introduces transactions that aggregate value in one, or few,
  coins, where the sender pays the receiver and gives 'change' back to
  themself. [Model1/README.md](Model1/README.md) explains how the first upgrade
  works.
- Model 2: extends the ledger to include blocks which are backchained
  by hash pointers, with each block containing a Merkle tree having
  2<sup>N</sup> transactions. [Model2/README.md](Model2/README.md)
  explains how the second upgrade works.
- *Model 3 (forthcoming): introduces payment by Scripts, where the new
  owner of a coin value stipulates the procedure by which it can be
  redeemed.*

Install, Compile, Explore:
-------------------------

The project is developed using [Go][go] version 1.6. Once you have that installed,
with your `$GOPATH` set, everything can be compiled using
```bash
$ python gobuilder.py
```

Individual Golang module functions are illuminated by:
```bash
$ python helpers.py unblock.go (or genesis.go, or ...)
```

[btc]: https://github.com/btcsuite/btcd
[go]: https://golang.org/

