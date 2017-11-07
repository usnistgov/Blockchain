Model 2.0 (April 21, 2017)
--------------------------

Bitcoin is a cryptocurrency whose principal recording mechanism is a ledger,
but a ledger organized as a list with $$\mathcal{O}(N)$$ access time tends
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
is '007' here. This takes typically between 100,000 and 1,000,000
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





Model 1.0 March 19, 2017
------------------------

Coin creation uses the "Trillion Dollar Platinum Coin" concept of
creating a single large value coin (here, 1 million units), and
spending it into existence incrementally through regular UBI
exercises. In version 0, CreateCoins used unit coins, and 1 million
outputs of denomination 1 require 1 million discrete outputs. Here,
the million unit coin is accomplished in a single line of output.

```
[471] Ledger:
    Tid: 1487547740451
    Ttyp: CreateCoins
    Id/Seq: 0/0 Denom: 0 Owner: MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQ
    CrtJSkjl+P06RUIACreJe7RTL5kBKUbYR+KNrvCbmjyKhvHbSWtLNAqrIHKyrcL0re
    eXfFS6IzFSxHsbr00b1t315IFxVuVADfnavZicTezMFyEapCDPDvJLxXeSPoXJ6hyx
    OB8SgeJ9Ted3HhNss4iUjs1pC1W85KKPuePFK4owIDAQAB
    Id/Seq: 1487547740451/0 Denom: 1000000 Owner: QSoT0u9VlnrL4wq2pppy
    +jB4lEbYJ7xeWnE1VKzgVic=
    Tsig:  N7zvw+ZyGzPlVCRy5cufve+UVhm5K+5PIn+zBH64iNfc7C4V9G9yYXiCLav
    fSmsCKKxFDspRGibolMvu8SwVyg5Bml5+ZaNlpBejgzJSNDcHd7kFuxUarOIQE5pWo
    6DJE1cHpPHcZdm0gafEdNXUWs3hJ3Ug7wkxVqKnqf4GBcY=
```

Large denomination coins achieve impressive orders of transaction size
compression, since no transaction requires more than 1 input and 2
outputs. This compares favorably against version 0, which for the
above requires 1 million inputs and 1 million outputs. This explains
why, not only is the domestic currency divided into larger and smaller
bills and coins, but also the checking system allows a single voucher
to transact arbitrarily large amounts of money.

Corollary to this improvement is that most transactions will now 'give
change' back to the payer. Much as you get $19 change back when you
give a $20 bill for a $1 candy bar, a payer who has a 20 unit coin to
hand over in payment for a 1 unit obligation, puts the 20 unit coin
as a single input, and generates 2 outputs: a 1 unit coin to the payee,
and a 19 unit coin back to self. Inputs and Outputs must balance, and
the chain of ownership is traceable, since the coin owner is identified 
by public key and the transaction outputs are signed by the corresponding
private key. Here is an instance of the UBI payment with the banker starting
at a 1010000 coin, distributing 20 to each recipient and returning
1009840 change to self.

```
[4] Ledger:
    Tid: 1487946650275
    Ttyp: PayCoins

Id/Seq: 14879270/0 Denom: 1010000 Owner: QSoT0

Id/Seq: 14879275/0 Denom: 20 Owner: wVWP7p
Id/Seq: 14879275/1 Denom: 20 Owner: XLuz8U
Id/Seq: 14879275/2 Denom: 20 Owner: 2MwzAJ
Id/Seq: 14879275/3 Denom: 20 Owner: /4fvUk
Id/Seq: 14879275/4 Denom: 20 Owner: tGjSEm
Id/Seq: 14879275/5 Denom: 20 Owner: vK1HKs
Id/Seq: 14879275/6 Denom: 20 Owner: uQREkQ
Id/Seq: 14879275/7 Denom: 20 Owner: jlvgVr
Id/Seq: 14879275/8 Denom: 1009840 Owner: QSoT0u

    Tsig:  VWrR5V ... g6GUQNStO5KiEcq5NYfTTZmc=
```

Components of the same remain the same as for Model 0.

Model 0, Feb 19, 2017
---------------------

Initial versions of the client, listener and blockchain clearing 
process.

- coin.go: is the client that can initiate Balance inquiries 
  PayCoins payments, Transactions inquiries, CreateCoins transactions,
  and remote listener closure (Quit). These commands can be initiated
  discretely from the command line, for any known payer/payee pair;
  in batch through a Multitest command listing specific PayCoins
  transactions one per line; or as a randomized automatic process
  through an Autotest command that generates an endless sequence
  of pairings with a random amount, within range of the payer's
  balance.

  PayCoins:  commands of the form "PayCoins alice.conf 3 chaz.pub"
  ./coin creates a transaction of type PayCoins, with alice as payer
  (the Inputs), chaz as payee (the Outputs) and for an amount of 3 
  coins. In this initial version, there is one coin of denomination 1,
  per input. In order to pay 3 coins, alice inquires her balance, and 
  if it exceeds 3, takes 3 coins and signs them over to chaz.

  Balance: commands of the form "Balance alice.conf", ./coin creates
  a transaction of type Balance, with zero inputs and zero outputs,
  and signs it.

  Transactions: commands of the form "Transactions chaz.conf", ./coin
  creates a transaction of type Transactions, with zero inputs and  zero
  outputs, and signs it. This yields a list of all transactions the
  user participates in, as Input giver (debits) or Output receiver 
  (credits).

  CreateCoins: commands of the form "CreateCoins 20 1 scrooge.conf",
  ./coin creates a transaction of type CreateCoins, with zero inputs and
  20 outputs of denomination 1 each, and signs it. The signature must be
  that of scrooge the banker.

  Quit: commands of the form "Quit scrooge.conf", ./coin creates a transaction
  of type Quit, with zero inputs and outputs, and signs it. The signature
  must be that of scrooge, the banker, the only party empowered to shut
  down the listener.

  Multitest: commands of the form "Multitest oneround.cli", ./coin parses
  the file of commands (PayCoins, CreateCoins, Quit, etc), and executes
  them sequentially.

  Autotest: commands of the form "Autotest users", ./coin constructs a 
  list of PayCoins commands with a round-robin of all combinations of
  the known users as payer and payee, with random payment values within
  their respective balances, and executes them sequentially. The
  randomized round-robin is repeated ad infinitum.

- masterpay.go: is the payments listener process that receives the
  set of transactions sent by ./coin, validates each transaction
  by checking the signature, determines if a PayCoins is within
  balance of the sender, and stores the received transactions in the
  newtxs new transactions file, as processed periodically by Clearing.

- The Clearing suite maintains the ledger, and comprises: genesis.go, 
  createcoins.go, ubi.go, paycoins.go.
  
  - The ledger does not exist and cannot grow without an initial Genesis
  transaction. genesis.go creates that. This will provide a termination
  point to the blockchain, once backward chaining hash pointers are in
  effect.

  - Coinage is created by the central banker with specific CreateCoins
  commands. In this version, all coins are of denomination 1, and
  payments are made by listing together a set of coins owned by the
  intending sender. The banker generates a sufficient set of coins via
  CreateCoins, for later distribution.

  - Coins are distributed equally from the stock held by the banker, among
  all known recipients, with the UBI or Universal Basic Income program.
  Input and Output coin quantities are equally balanced in a PayCoins
  transaction, signed over to the new recipient by the banker. The
  remainder left over after distribution is retained by the banker, until
  after the next CreateCoins and UBI command executions.

  - The paycoins process is run as a cronjob to periodically process new
  transactions approved by masterpay, the listener. PayCoins transactions
  are approved (mostly: insufficient funds and double spend attempts
  are filtered by masterpay), and added to the ledger. In this first
  instance, the ledger is implemented as a forward chained list of
  transactions. Organization of transactions as a Merkle Tree, with 
  backward chained hash pointers implementing a blockchain, will be
  realized in a subsequent version.


- Ancillary tools include: listall.go, sum.go, balances.go, mycoins.go

  - listall.go <ledger or newtxes>: do a full print of the fields in
    every transaction in the ledger, or every new transaction to be
    processed.

  - sum.go <ledger or newtxes>: Summarize each transaction on one line,
    with Id/Seq, transaction name, verification status, #of inputs,
    #of outputs, short pub key Id.

  - balances.go: list the balance of all accounts in the system,
    including the banker's (undistributed) balance.

  - mycoins.go <my public key file>: list the balance of the account
    idetified by public key file.


- Running the System:

  The listener masterpay.go runs as a background process, listening for
  new transaction requests from the (or any) client, acts on them, replies
  to the client, and saves the resulting transactions in a newtxs file.
  The client runs on demand from a user to generate those requests. This
  alternating process between coin client and masterpay listener serves
  to set up balance changes among the members, but newtxs must be
  serviced regularly by a clearing process. Create Coins and UBI
  distribution episodically increase the total coin balance available
  for the members to transact with. Paycoins is the heart of the clearing
  process which runs sufficiently often to satisfy of clear' the new
  transactions, adjust the members' balances and allow the members to
  keep paying and receiving among themselves.



- The structure of Each Transaction: PayCoins, CreateCoins, Genesis [0]

[0] Ledger:
    Tid: 0
    Ttyp: BigBang
    Id/Seq: 0/0 Denom: 0 Owner: MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBg
    QCrtJSkjl+P06RUIACreJe7RTL5kBKUbYR+KNrvCbmjyKhvHbSWtLNAqrIHKyrcL0
    reeXfFS6IzFSxHsbr00b1t315IFxVuVADfnavZicTezMFyEapCDPDvJLxXeSPoXJ6
    hyxOB8SgeJ9Ted3HhNss4iUjs1pC1W85KKPuePFK4owIDAQAB
    Id/Seq: 0/0 Denom: 0 Owner: v=HoR
    Tsig:  jjIOdnyx8v3L365zuMzi3zvvjDfroClpHnljtPOZCk0Dnsrj0AH1KBSlNB
    vAzo97Hrs1PZexLbp1Fk7NK2MewitsIwkXXdPSOndRdfqIiw7u4zPTzo6QozP5/gh
    5O6KiXn97GhtSUYn9GcVUrA6koD9y7MBesj0OZpkaBRaXki8=

[1] Ledger:
    Tid: 1487272161857
    Ttyp: CreateCoins
    Id/Seq: 0/0 Denom: 0 Owner: MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBg
    QCrtJSkjl+P06RUIACreJe7RTL5kBKUbYR+KNrvCbmjyKhvHbSWtLNAqrIHKyrcL0
    reeXfFS6IzFSxHsbr00b1t315IFxVuVADfnavZicTezMFyEapCDPDvJLxXeSPoXJ6
    hyxOB8SgeJ9Ted3HhNss4iUjs1pC1W85KKPuePFK4owIDAQAB
    Id/Seq: 1487272161857/0 Denom: 1 Owner: QSoT0u9VlnrL4wq2pppy+jB4l
    EbYJ7xeWnE1VKzgVic=
    Id/Seq: 1487272161857/1 Denom: 1 Owner: QSoT0u9VlnrL4wq2pppy+jB4l
    EbYJ7xeWnE1VKzgVic=
     . . .
    Id/Seq: 1487272161857/48 Denom: 1 Owner: QSoT0u9VlnrL4wq2pppy+jB4l
    EbYJ7xeWnE1VKzgVic=
    Id/Seq: 1487272161857/49 Denom: 1 Owner: QSoT0u9VlnrL4wq2pppy+jB4l
    EbYJ7xeWnE1VKzgVic=
    Tsig:  Sc9nz3h+QkQgwjQ2rYJTW56egzJkZXp2B2r1imEvFSNyNRVkxh4nR7CbBnN
    VDX2aj4McQEHb6JDo+y5wLCpzog7k7Lb0MbnxSPWmGVRrprQcKpDIX72/2Rnis7N87
    E01gyVePK5Nd6BiLOkvzBlKHCSCAulgf2HRQh8rYiTxwuM=

[32] Ledger:
    Tid: 1487362712339
    Ttyp: PayCoins
    Id/Seq: 1487272161857/8 Denom: 1 Owner: MIGfMA0GCSqGSIb3DQEBAQUAA4
    GNADCBiQKBgQChw+/3t6K8N30eKEbYaE0ptaDqOEONobsQZ+blpwl7ADPgqS+qxi6p
    Xt1meR8WybK/I71OZM6fhqgjr5dLK1+GF8fkeaeBbU/YGYVP2ChJ1dp/Ju7ku3KfgE
    YyL48d2AuwDSAyxkYyffwbGIiJg3Li+SPBvRLUigazPvPxCtlCMwIDAQAB
    Id/Seq: 1487272161857/8 Denom: 1 Owner: tGjSEmtsm4D8A4M38GwE2j7Rwm
    LyoIXzpLdOmYNQSvs=
    Tsig:  n0zr+pQWo4gCkpMnJd77E/aq3OpPQlTL/5ToDlnmHRTteo4Ujx5T9yEjCEW
    loDaNyiTAjeKJRORPrvtB6IL4TiaDVd4JdfM2i6Gn6+3QT0+IdDsLT6cxFOlJaUfSL
    weu70yjOimehNqijLh1i02bEZOJacblCvGMe2jFKynpFoA=


- Additional client-listener 'protocol' transactions: Quit, Balance,
  Transactions, Error, PayeeReceipt, CreateReceipt : 
  Quit, Balance, Transactions are created in the format of a Currency
  system transaction, with Inputs, Outputs and a Signature.
  Error, PayeeReceipt and CreateReceipt are replies from the Listener
  to the Client, also in the format of a Currency transaction.



