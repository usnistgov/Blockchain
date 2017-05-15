/*
************************************************************************
from currency.go
  ADD: makegenesis.go
       Test making genesis.
       Verify it including the signature and output value.
       (func FindGenesis(str.Tr) bool)
  genesis.go:
    create gensis if not exists.
    load and store ledger; print genesis only.
    genesis integrity tests: Ttyp == "BigBang"; outputs[0] ==" v=HoR"

  ADD: makecreate.go
       Test creating coins.
       Verify the trans including signature and total created.
       Find the total number of coins in the system.
  ADD: makedistribute.go
       For all known users, create a block of coins.
       Confirm total of coins distributed + unassigned = total created.
       Confirm transfer of ownership from banker to members.
       Initial test is to distribute one coin per member,
       and leave one remainder for banker.
  add: makepay.go
       Create a PayCoins of one coin transferred. 
       Confirm before and after balances of payor and payee.
       Create a PayCoins of two coins transferred.
       Confirm before and after.
       Create and confirm PayCoins in increments up to 6.
       Create multiple PayCoins transactions in increments to 6,
       with varying amounts. Confirm before and after balances.

  methods.GetConfigs
       Check values for multiple config files.
       Consider including function return values for PubKey operations.
  structures.PrintMap
       Print a map[string]string.
       Can this method be generalized?
  methods.CheckErrorInst
       Test with nil and positive error values.
  structures.PrintLedger
       Print ledgers of length 1-N.
       Can the print be compacted without loss of information?
  structures.PrintCoins
       Print coin structures of length 1-N.
  newtxes.CreateGenesis
       Check ledger existence before and after.
       Verify signature and Output value.
  newtxes.CreateCoins
       As per makecreate above.
  methods.HashSigs
  methods.CoinBase
       Group Coinbase operations together:
       struct.Coin,
       CreateCoins,
       PrintCoins,
       PayCoins,
       DestroyCoins (i.e. a tax)?
       CountCoins,
       GetMemberBalances.
  newtxes.DistributeUBI
  methods.ProcessNewTransactions
       Should be constructed from already tested elements,
       and involves inbound transactions to add to the ledger.
  methods.StoreLedger
       Group ledger operations together:
       struct.Transaction,
       CreateLedger,
       LoadLedger,
       PrintLedger,
       AddToLedger,
       StoreLedger.
************************************************************************

balances.go
  structures.Transaction
  structures.Coin
  methods.LoadLedger
  methods.M1

************************************************************************
coinbase.go
  <balances>
  structures.PrintCoins

************************************************************************
listall.go
  structures.ShortTransaction
  structures.PrintTransaction

************************************************************************
mycoins.go
 <balances>
  himitsu.HashPublicKey
  methods.GetMyCoins
  structures.CoinCount

************************************************************************
sum.go
  Additional cross checks needed in VerifyTransaction:
  ADD: VerifyGenesis
  ADD: VerifyCreate
  ADD: VerifyPay

  methods.VerifyTransaction

************************************************************************
coin.go
  Test wrapping and unwrapping, sending, receiving and verifying txs.
  Send multiple txes per connection, multiple connections.
  methods.GetDir
  structures.CoinValues
  methods.NoFile
  newtxes.GetQuit
  newtxes.GetBalance
  himitsu.BaseDER
  structures.PrintCoin
  newtxes.PayCoins
  methods.CheckError

************************************************************************
masterpay.go
  structures.PrintMap
  methods.M1
  methods.GetMyCoins
  methods.DoSign
************************************************************************
*/

