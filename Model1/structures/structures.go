package structures

import( 
  "container/list"
  "fmt"
  "strconv"
) // end import.


// good for all Transactions: Genesis, CreateCoins, PayCoins.
type Transaction struct {
  Tid int64
  Ttyp string
  Inputs []Coin
  Outputs []Coin
  Tsig string  // base64, for now.
} // end Transaction.

// Used as a slice in Inputs, Outputs.
type Coin struct {
  Cid int64
  Seq string
  Denom string
  Owner string
} // type Coin.


// Print the contents of a slice of Coin structs.
func Printslice(x []Coin, title string) {
  fmt.Println("\n[", title)
  for ix := 0; ix < len(x); ix++ {
    fmt.Println("\t", x[ix])
  }
  fmt.Println("]")
} // end Printslice.


// Print a single Transaction of the above types.
func PrintTransaction(gx Transaction, sx string) {
  fmt.Println(sx)
  fmt.Println("    Tid:", gx.Tid)
  fmt.Println("    Ttyp:", gx.Ttyp)
  fmt.Println("")
  PrintCoins(gx.Inputs, "OwnerIn ")
  fmt.Println("")
  PrintCoins(gx.Outputs, "OwnerOut ")
  fmt.Println("    Tsig: ", gx.Tsig)
} // printTransaction.

// Print a Short form Transaction.
func ShortTransaction(gx Transaction) {
  fmt.Printf("Tid: %d, %s, Outputs:\n", gx.Tid, gx.Ttyp)
  PrintCoins(gx.Outputs, "OwnerOut")
} // ShortTransaction.


// Print the entire ledger of Transactions.
func PrintLedger(led []Transaction) {
  for ix := 0; ix < len(led); ix++ {
    PrintTransaction(led[ix], led[ix].Ttyp)

  } // end for.

} // end PrintLedger.

// Print a generalized List.
func Printlist(x *list.List, title string) {
  fmt.Println("\n[", title)
  for e := x.Front(); e != nil; e = e.Next() {
    fmt.Println("\t", e.Value)
  }
  fmt.Println("]")
} // end Printlist.

// Print my coins.
func PrintMyCoins(pcoins []Coin, mypub string) {
  for px := 0; px < len(pcoins); px++ {
    if pcoins[px].Owner == mypub {
      PrintCoin(pcoins[px], "MyCoin")
    } // end if mypub.
  } // end for px.
} // end func PrintCoins.

// Print all coins in a coin transaction.
func PrintCoins(pcoins []Coin, InOrOut string) {
  fmt.Printf("[ %s\n", InOrOut)
  for px := 0; px < len(pcoins); px++ {
    PrintCoin(pcoins[px], InOrOut)
  } // end for px.
  fmt.Println(" ]")
} // end func PrintCoins.

// Print the Id, denomination and owner of a single coin.
func PrintCoin(thecoin Coin, InOrOut string) {
  fmt.Printf("Id/Seq: %013d/%s Denom: %s %s: %s\n", thecoin.Cid, thecoin.Seq, thecoin.Denom, InOrOut, thecoin.Owner)
} // end func PrintCoin.


// Return the count of coins in a structure.
func CoinCount(pcoins []Coin) int {
  pkval := 0
    
  for pk := 0; pk < len(pcoins); pk++ {
    pcval, _ := strconv.Atoi(pcoins[pk].Denom)
    pkval = pkval + pcval
  } // end for pk.

  return pkval

} // end func CoinCount.

// Return the total value of coins in a structure.
func CoinValues(pcoins []Coin) int {
  pcount := 0

  if len(pcoins) == 0 {
    return 0
  } //endif no coins.

  for px := 0; px < len(pcoins); px++ {
    thecoin := pcoins[px]
    coinvalue, _ := strconv.Atoi(thecoin.Denom)
    pcount += coinvalue
  } // end for pcoins.

  return pcount

} // end func CoinValues.


//Print the Keys and Values of the Map, and its Title.
func PrintMap(thismap map[string]string, thistitle string) {
  fmt.Println(thistitle)

  for key, value := range thismap {
    fmt.Println("\t", key, "=", value)
  } // for thismap.

} // end func PrintMap.

