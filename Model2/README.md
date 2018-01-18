# Model 2.0

Bitcoin is a cryptocurrency whose principal recording mechanism is a ledger,
but a ledger organized as a list with `O(N)` access time tends
to be an inefficient structure. The big innovation was to organize the ledger
as a blockchain where each succeeding block has a hash pointer, pointing back
to the previous block (the pointer is the hash of the previous block).
Transactions that have been verified by the Banker are organized into a
structure called a [Merkle tree][merkle]: every transaction is hashed, and
these transactions form the leaves of the tree. Pairs of hashes are grouped
together and hashed again, to form the next higher level of the tree. The
single, top level hash is the Merkle root, which is stored as a field of
the Block Header. If any transaction were to be changed, every hash in that
branch of the tree,up to the root, must be recomputed. Every block's
hash pointer, up to the end of the block chain, must be recomputed, in
order to preserve the integrity of the ledger. Model 2 in this series
implements the blockchain with backward chaining hash pointers, and
transactions organized as Merkle trees.

The implementation of Merkle trees in this model is perhaps non-standard
in that each one requires an exact power of two transactions. The hash
tree is constructed in levels, from the leaves up to the root. The entire
tree is stored in a map, allowing direct access to a transaction through
knowing its hash. The Block header and its hash are also stored in the
map. The genesis block and its succeeding block header are illustrated
below. Note the last block hash of the genesis block is null, as there
is no previuos block. The last block hash of the next block is the
hash of the genesis block. Proof-of-work is performed in this
implementation by matching the hash against a given prefix, which
is `007` here. This takes typically between 100,000 and 1,000,000
hash attempts to find, incrementing a nonce each time.

```
ThisBlockHash: 0078tnuxo/73t2a23eORhWEYK5Ia2O7uiKEN+CS2+Fs=
Block: 1492759091099
	Version: 1, BlockTime: 1492759091099, Nbits: 007, Nonce: 103950
	LastBlockHash: 0
	MerkleHash: zpgFzEMF1u/i4uJvfMrS1a/quo51lO1F6QMAgacB0Ro=
```

```
ThisBlockHash: 007U+xWZylTpz5/r/KU809k4Rb1VDk6jNfda7Pgt+0Y=
Block: 1492759092538
	Version: 1, BlockTime: 1492759092538, Nbits: 007, Nonce: 379710
	LastBlockHash: 0078tnuxo/73t2a23eORhWEYK5Ia2O7uiKEN+CS2+Fs=
	MerkleHash: pSRHTig1Ew81nq4h0rceUQrUxXoGRtrpjp3O+XyHy5Y=
```

In previous models, the ledger was stored in a file as a list of transactions,
each encoded in base64. In this model each block is stored as a list of
newline terminated base64 objects where each object may be either a hash,
a Block header, a concatenated pair of Merkle level hashes, or a
transaction. They unpack into a map which can be traversed given the
original toplevel block hash. The Merkle hash is retrieved from the Block
header, and the Merkle tree can also be traversed from there.

There are extensions and rewrites of the supporting functions, briefly
noted here.

- `./coinbase platinum.conf`: coinbase is the renamed createcoins. It enables
the banker to create coin value, for distribution to other participants.
- `./payments paycoins.conf`: the renamed paycoins. It clears new transactions,
adds them to the ledger, and forms blocks when 2**N transactions are
received.
- `./newbalances bals.conf`: reports the current balance of every participant,
based on the last block utxos and the current residual ledger.
- `./unblock blocks block [verbose|levels]`: ripple through the blockchain
unpacking and printing the block headers, which include the previous
block pointer, and the Merkle Root pointer. The verbose option unpacks
the Merkle Tree in each block and prints every transaction. The levels
option unpacks and prints the hashes at every level in the Merkle Tree
hierarchy.
- `./unblockutxos bals.conf`: using the last block utxos and the residual ledger,
unblockutxos finds and prints all current, unspent transactions and their
owners.

[merkle]: https://en.wikipedia.org/wiki/Merkle_tree
