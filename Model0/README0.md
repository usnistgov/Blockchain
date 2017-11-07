Model 0.0 (February 19, 2017)
-----------------------------

Initial versions of the client, listener and blockchain clearing
process.

- `coin.go` is the client that can initiate Balance inquiries
  PayCoins payments, Transactions inquiries, CreateCoins transactions,
  and remote listener closure (Quit). These commands can be initiated
  discretely from the command line, for any known payer/payee pair;
  in batch through a Multitest command listing specific PayCoins
  transactions one per line; or as a randomized automatic process
  through an Autotest command that generates an endless sequence
  of pairings with a random amount, within range of the payer's
  balance.

  - PayCoins:  commands of the form
    `"PayCoins alice.conf 3 chaz.pub", ./coin`
    creates a transaction of type PayCoins, with alice as payer
    (the Inputs), chaz as payee (the Outputs) and for an amount of 3
    coins. In this initial version, there is one coin of denomination 1,
    per input. In order to pay 3 coins, alice inquires her balance, and
    if it exceeds 3, takes 3 coins and signs them over to chaz.

  - Balance: commands of the form `"Balance alice.conf", ./coin`
    creates a transaction of type Balance, with zero inputs and zero outputs,
    and signs it.

  - Transactions: commands of the form `"Transactions chaz.conf", ./coin`
    creates a transaction of type Transactions, with zero inputs and
    zero outputs, and signs it. This yields a list of all transactions
    the user participates in, as Input giver (debits) or Output receiver
    (credits).

  - CreateCoins: commands of the form `"CreateCoins 20 1 scrooge.conf", ./coin`
    creates a transaction of type CreateCoins, with zero inputs and
    20 outputs of denomination 1 each, and signs it. The signature must be
    that of scrooge the banker.

  - Quit: commands of the form `"Quit scrooge.conf", ./coin` creates a
    transaction of type Quit, with zero inputs and outputs, and signs it.
    The signature must be that of scrooge, the banker, the only party
    empowered to shut down the listener.

  - Multitest: commands of the form `"Multitest oneround.cli", ./coin` parses
    the file of commands (PayCoins, CreateCoins, Quit, etc), and executes
    them sequentially.

  - Autotest: commands of the form `"Autotest users", ./coin` constructs a
    list of PayCoins commands with a round-robin of all combinations of
    the known users as payer and payee, with random payment values within
    their respective balances, and executes them sequentially. The
    randomized round-robin is repeated ad infinitum.

- `masterpay.go` is the payments listener process that receives the
  set of transactions sent by `./coin`, validates each transaction
  by checking the signature, determines if a PayCoins is within
  balance of the sender, and stores the received transactions in the
  newtxs new transactions file, as processed periodically by Clearing.

- The Clearing suite maintains the ledger, and comprises: `genesis.go`,
  `createcoins.go`, `ubi.go`, `paycoins.go`.

  - The ledger does not exist and cannot grow without an initial Genesis
    transaction. `genesis.go` creates that. This will provide a termination
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


- Ancillary tools include: `listall.go`,
                           `sum.go`,
                           `balances.go`,
                           `mycoins.go`

  - `listall.go <ledger or newtxes>`: do a full print of the fields in
    every transaction in the ledger, or every new transaction to be
    processed.

  - `sum.go <ledger or newtxes>`: summarize each transaction on one line,
    with Id/Seq, transaction name, verification status, #of inputs,
    #of outputs, short pub key Id.

  - `balances.go`: list the balance of all accounts in the system,
    including the banker's (undistributed) balance.

  - `mycoins.go <my public key file>`: list the balance of the account
    idetified by public key file.


- Running the System:

  The listener `masterpay.go` runs as a background process, listening for
  new transaction requests from the (or any) client, acts on them, replies
  to the client, and saves the resulting transactions in a `newtxs` file.
  The client runs on demand from a user to generate those requests. This
  alternating process between coin client and masterpay listener serves
  to set up balance changes among the members, but newtxs must be
  serviced regularly by a clearing process. Create Coins and UBI
  distribution episodically increase the total coin balance available
  for the members to transact with. Paycoins is the heart of the clearing
  process which runs sufficiently often to satisfy or clear the new
  transactions, adjust the members' balances and allow the members to
  keep paying and receiving among themselves.

- The structure of Each Transaction: PayCoins, CreateCoins, Genesis

```
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
```

```
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
```

```
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
```

- Additional client-listener 'protocol' transactions:
  - Quit, Balance, Transactions are created in the format of a Currency
  system transaction, with Inputs, Outputs and a Signature.
  - Error, PayeeReceipt and CreateReceipt are replies from the Listener
  to the Client, also in the format of a Currency transaction.
