package newtxes

import (
  "fmt"
  "currency/himitsu"
  "currency/methods"
  "currency/structures"
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

    
// PayCoins transaction. Transfer from old owner to new owner.
// Old owner signs the outputs.
// New owner is different from old owner.
func PayCoins(conmap map[string]string, coinids []structures.Coin, newpub string) structures.Transaction {

  // Deep copy coinids to paidcoins + new owner.
  paidcoins := []structures.Coin{}
  for px := 0; px < len(coinids); px++ {
    oldone := coinids[px]
    tfrd := structures.Coin{oldone.Cid, oldone.Seq, oldone.Denom, newpub}
    paidcoins = append(paidcoins, tfrd)
  } // end for paidcoins.
  sig := methods.DoSign(conmap, paidcoins)
  transid := methods.MilliNow()
  tx := structures.Transaction{transid, "PayCoins", coinids, paidcoins, sig}
  return tx

} // end PayCoins.

// Distribute the inventory of unclaimed coins among the known beneficiaries.
// Any unclaimed coins are kept by the banker.
func DistributeUBI(conmap map[string]string, supply []structures.Coin, bennies []string) []structures.Transaction {
  newbundle := []structures.Transaction{}

  if len(bennies) > len(supply) {
    fmt.Println("Not enough coins to distribute.")
    return newbundle   // length 0.
    } // endif no supply.

  quot, _ := DivRem(len(supply), len(bennies))
  for bx := 0; bx < len(bennies); bx++ {
    payout := Disburse(conmap, supply[(bx*quot):], quot)
    nextrx := PayCoins(conmap, payout, bennies[bx])
    newbundle = append(newbundle, nextrx)
  } // end for bennies.

  return newbundle

} // end func DistributeUBI.


// Return the quotient and remainder.
func DivRem(sor int, dend int) (int, int) {
  quot := sor / dend
  rem := sor % dend

  return quot, rem

} // end func DivRem.


// Distribute the quotient of coins.
// Change Owner to long form public key.
func Disburse(cmap map[string]string, sup []structures.Coin, qt int) []structures.Coin {

  ubipay := []structures.Coin{}

  for ux := 0; ux < qt; ux++ {
    sup[ux].Owner = himitsu.BaseDER(cmap["pubkey"])
    ubipay = append(ubipay, sup[ux])
  } // end quot.

  return ubipay

} // end func Disburse.


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

