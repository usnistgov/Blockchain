package newtxes

import (
  "fmt"
  "currency/himitsu"
  "currency/methods"
  "currency/structures"
  "os"
  "strconv"
)



// Create the Genesis transaction.
func CreateGenesis(conmap map[string]string) structures.Transaction {
  owner := himitsu.BaseDER(conmap["pubkey"])
  thinair := structures.Coin{ 0, "0", "0", owner }
  inputs := []structures.Coin{ thinair }
  thickair := structures.Coin{ 0, "0", "0", "v=HoR" }
  outputs := []structures.Coin{ thickair }
  sig := methods.DoSign(conmap, outputs)
  tx := structures.Transaction{0, "BigBang", inputs, outputs, sig}
  return tx
} // end createGenesis.



// CreateCoins transaction. Create Coins up to coin * denom args.
// Assign them to Scrooge as owner.
// So new owner is the same as old owner.
func CreateCoins(conmap map[string]string) structures.Transaction {
  owner := himitsu.BaseDER(conmap["pubkey"])
  ownhash := himitsu.HashPublicKey(conmap["pubkey"])
  thinair := structures.Coin{ 0, "0", "0", owner }
  inputs := []structures.Coin{ thinair }
  thickair, transid := MakeCoins(conmap, ownhash)
  sig := methods.DoSign(conmap, thickair)
  tx := structures.Transaction{transid, "CreateCoins", inputs, thickair, sig}
  return tx
} //end CreateCoins.

// MakeCoins: Get coincount and denomination from config file and generate
// one output for each new coin. Assign them to Scrooge's public key.
func MakeCoins(conmap map[string]string, owner string) ([]structures.Coin, int64) {
  coins, _ := strconv.Atoi(conmap["coin"])
  trid := methods.MilliNow()
  denom := conmap["denom"]
  newcoins := []structures.Coin{}

  for seq := 0; seq < coins; seq++ {
    stseq := strconv.Itoa(seq)
    onecoin := structures.Coin{trid, stseq, denom, owner} 
    newcoins = append(newcoins, onecoin)
  } // end for seq.

  return newcoins, trid

} // end MakeCoins.

    

// DistributeUBI: disburse a quotient of coins to each listed
// beneficiary, from the banker's supply. Consolidate into one
// coin if necessary.
func DistributeUBI(conmap map[string]string, supply []structures.Coin, bennies []string, ubiquot int) []structures.Transaction {

  newbundle := []structures.Transaction{}
  transout := structures.Transaction{}
  oneout := structures.Coin{}

  // Is there enough value to give out?
  if len(bennies) * ubiquot > structures.CoinValues(supply) {
    fmt.Println("The Cupboard is bare.")
    os.Exit(1)
  } // endif nosupply.

  if len(supply) == 1 {
    oneout = supply[0]
    oneout.Owner = himitsu.BaseDER(conmap["pubkey"])
  } else {
    for sx := 0; sx < len(supply); sx++ {
      denom, _ := strconv.Atoi(supply[sx].Denom)
      if denom > len(bennies) * ubiquot {
        oneout = supply[sx]
        oneout.Owner = himitsu.BaseDER(conmap["pubkey"])
      } // endif denom.
    } // end for supply.
  } // endif supply.

  if  oneout == (structures.Coin{}) {
    oneout, transout = Consolidate(conmap, supply)
    newbundle = append(newbundle, transout)
  } // endif nil.

  // carve out one UBI payout per beneficiary.
  nextrx := PayOut(conmap, oneout, ubiquot, bennies)
  newbundle = append(newbundle, nextrx) 
  
  return newbundle

} // end func DistributeUBI.

// Consolidate: if owner has more than one coin, combine all the
// values into a single coin. Write out as a new transaction.

func Consolidate(cmap map[string]string, supplyin []structures.Coin) (structures.Coin, structures.Transaction) {

  thistime := methods.MilliNow()
  owner := himitsu.HashPublicKey(cmap["pubkey"])
  aggregate := 0

  // aggregate the input values into a single output value.
  for sx := 0; sx < len(supplyin); sx++ {
    value, _ := strconv.Atoi(supplyin[sx].Denom)
    aggregate += value
    // while we're atit let's make sure the input has full pubkey.
    supplyin[sx].Owner = himitsu.BaseDER(cmap["pubkey"])
  } // end for supplyin.

  aggy := fmt.Sprintf("%d", aggregate)
  consol := structures.Coin{thistime, "0", aggy, owner}
  outputs := []structures.Coin{ consol }
  tsig := methods.DoSign(cmap, outputs)
  constrx := structures.Transaction{thistime, "PayCoins", supplyin, outputs, tsig}
  structures.PrintTransaction(constrx, "Consolidated Coinbase:")

  return consol, constrx

} // end func Consolidate.


// PayOut takes a single large coin and makes payments to all beneficiaries
// in one transaction.
func PayOut(conmap map[string]string, onein structures.Coin, payment int, payees []string) structures.Transaction {

  // Owner DER?
  inputs := []structures.Coin{}
  inputs = append(inputs, onein)

  // create the outputs vector, inc. change.
  outputs := []structures.Coin{}
  now := methods.MilliNow()
  valuein, _ := strconv.Atoi(onein.Denom)
  cseq := 0
  ascpay := fmt.Sprintf("%d", payment)
  
  // hand out the payment to all payees.
  for px := 0; px < len(payees); px++ {
    paycoin := structures.Coin{now, fmt.Sprintf("%d", cseq), ascpay, payees[px]}
    outputs = append(outputs, paycoin)
    cseq += 1
    valuein = valuein - payment
  } // end for payees.

  ownerpub := himitsu.HashPublicKey(conmap["pubkey"])
  changecoin := structures.Coin{now, fmt.Sprintf("%d", cseq), fmt.Sprintf("%d", valuein), ownerpub}
  outputs = append(outputs, changecoin)
  tsig := methods.DoSign(conmap, outputs)
  trx := structures.Transaction{now, "PayCoins", inputs, outputs, tsig}

  return trx

} // end func PayOut.


// PayCoins takes a single large coin, makes the payment to payee, and
// gives the balance back to the owner.
// payee needs to be the hashed form of the public key.

func PayCoins(conmap map[string]string, onein structures.Coin, payment int, payee string) structures.Transaction {

  // Create the inputs vector.
  inputs := []structures.Coin{}
  inputs = append(inputs, onein)

  // Create the outputs vector.
  outputs := []structures.Coin{}
  now := methods.MilliNow()
  valuein, _ := strconv.Atoi(onein.Denom)
  change := valuein - payment
  paycoin := structures.Coin{now, "0", fmt.Sprintf("%d", payment), payee}
  outputs = append(outputs, paycoin)
  ownerpub := himitsu.HashPublicKey(conmap["pubkey"])
  changecoin := structures.Coin{now, "1", fmt.Sprintf("%d", change), ownerpub}
  outputs = append(outputs, changecoin)

  // Sign the outputs and create the transaction.
  tsig := methods.DoSign(conmap, outputs)
  trx := structures.Transaction{now, "PayCoins", inputs, outputs, tsig}

  return trx

} //end func PayCoins.


// GetClientCommand: ancillary transactions, used in Client-Listener,
// Not used in Blockchain updater.
// Aggregates: Getbalance, GetQuit, GetTransactions

func GetClientCommand(conmap map[string]string, clientcmd string) structures.Transaction {

  owner := himitsu.HashPublicKey(conmap["pubkey"])
  ownair := himitsu.BaseDER(conmap["pubkey"])
  thinair := structures.Coin{0, "0", "0", ownair}
  inputs := []structures.Coin{ thinair }
  thickair := structures.Coin{0, "0", "0", owner}
  outputs := []structures.Coin{ thickair }
  sig := methods.DoSign(conmap, outputs)
  tx := structures.Transaction{0, clientcmd, inputs, outputs, sig}
  return tx

} // end func GetClientCommand.

// GetBalance: ancillary transactions, used in Client-Listener,
// Not used in BLockchain updater.
func GetBalance(conmap map[string]string) structures.Transaction {

  owner := himitsu.HashPublicKey(conmap["pubkey"])
  ownair := himitsu.BaseDER(conmap["pubkey"])
  thinair := structures.Coin{0, "0", "0", ownair}
  inputs := []structures.Coin{ thinair }
  thickair := structures.Coin{0, "0", "0", owner}
  outputs := []structures.Coin{ thickair }
  sig := methods.DoSign(conmap, outputs)
  tx := structures.Transaction{0, "Balance", inputs, outputs, sig}
  return tx

} // end func GetBalance.


// GetBalance: ancillary transactions, used in Client-Listener,
// Not used in BLockchain updater.
func GetQuit(conmap map[string]string) structures.Transaction {

  owner := himitsu.HashPublicKey(conmap["pubkey"])
  ownair := himitsu.BaseDER(conmap["pubkey"])
  thinair := structures.Coin{0, "0", "0", ownair}
  inputs := []structures.Coin{ thinair }
  thickair := structures.Coin{0, "0", "0", owner}
  outputs := []structures.Coin{ thickair }
  sig := methods.DoSign(conmap, outputs)
  tx := structures.Transaction{0, "Quit", inputs, outputs, sig}
  return tx

} // end func GetQuit.



// CloseListener: ancillary transactions, used in Client-Listener,
// Not used in Blockchain updater.
func CloseListener(conmap map[string]string) structures.Transaction {

  owner := himitsu.BaseDER(conmap["pubkey"])
  thinair := structures.Coin{0, "0", "0", owner}
  inputs := []structures.Coin{ thinair }
  thickair := structures.Coin{0, "0", "0", owner}
  outputs := []structures.Coin{ thickair }
  sig := methods.DoSign(conmap, outputs)
  tx := structures.Transaction{0, "Quit", inputs, outputs, sig}
  return tx

} // end func CloseListener.

