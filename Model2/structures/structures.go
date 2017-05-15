package structures

/*
This software was developed by employees of the National Institute of 
Standards and Technology (NIST), an agency of the Federal Government. 
Pursuant to title 17 United States Code Section 105, works of NIST 
employees are not subject to copyright protection in the United States 
and are considered to be in the public domain. Permission to freely 
use, copy, modify, and distribute this software and its documentation 
without fee is hereby granted, provided that this notice and disclaimer 
of warranty appears in all copies.

THE SOFTWARE IS PROVIDED 'AS IS' WITHOUT ANY WARRANTY OF ANY KIND, 
EITHER EXPRESSED, IMPLIED, OR STATUTORY, INCLUDING, BUT NOT LIMITED TO, 
ANY WARRANTY THAT THE SOFTWARE WILL CONFORM TO SPECIFICATIONS, ANY 
IMPLIED WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, 
AND FREEDOM FROM INFRINGEMENT, AND ANY WARRANTY THAT THE DOCUMENTATION 
WILL CONFORM TO THE SOFTWARE, OR ANY WARRANTY THAT THE SOFTWARE WILL BE 
ERROR FREE. IN NO EVENT SHALL NIST BE LIABLE FOR ANY DAMAGES, INCLUDING, 
BUT NOT LIMITED TO, DIRECT, INDIRECT, SPECIAL OR CONSEQUENTIAL DAMAGES, 
ARISING OUT OF, RESULTING FROM, OR IN ANY WAY CONNECTED WITH THIS SOFTWARE, 
WHETHER OR NOT BASED UPON WARRANTY, CONTRACT, TORT, OR OTHERWISE, WHETHER 
OR NOT INJURY WAS SUSTAINED BY PERSONS OR PROPERTY OR OTHERWISE, AND 
WHETHER OR NOT LOSS WAS SUSTAINED FROM, OR AROSE OUT OF THE RESULTS OF, 
OR USE OF, THE SOFTWARE OR SERVICES PROVIDED HEREUNDER.
*/

/*
  Stephen Nightingale
  night@nist.gov
  NIST, Information Technology Laboratory
  March 1, 2017
*/


import( 
  "container/list"
  "fmt"
  "strconv"
  "strings"
  "currency/himitsu"
) // end import.


// The primordial Block, length 104 bytes.
type Block struct {
  Version int        // 0000000001
  LastBlockHash string // len 44 bytes
  MerkleHash string    // len 44 bytes
  BlockTime int64      // Unix time
  Nbits string          // Target threshold
  Nonce int          // Fudger to achieve target threshold
} // end Block.



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


// Print the contents of a slice of strings.
func Stringslice(x []string, title string) {
  fmt.Println("\n[", title)
  for ix := 0; ix < len(x); ix++ {
    fmt.Println("\t", x[ix])
  }
  fmt.Println("]")
} // end Stringslice.

// Print the contents of a slice of Coin structs.
func Printslice(x []Coin, title string) {
  fmt.Println("\n[", title)
  for ix := 0; ix < len(x); ix++ {
    fmt.Println("\t", x[ix])
  }
  fmt.Println("]")
} // end Printslice.


/* Print a Block Header containing the following fields:
  Version int        // 0000000001
  LastBlockHash string // len 44 bytes
  MerkleHash string    // len 44 bytes
  BlockTime int64      // Unix time
  Nbits int          // Target threshold
  Nonce int          // Fudger to achieve target threshold
*/
func PrintBlockHeader(bx Block) {

  fmt.Printf("Block: %d\n",  bx.BlockTime)
  fmt.Printf("\tVersion: %d, BlockTime: %d, Nbits: %s, Nonce: %d\n",  bx.Version, bx.BlockTime, bx.Nbits, bx.Nonce)
  if bx.LastBlockHash == "" {
    fmt.Printf("\tLastBlockHash: 0\n")
  } else {
    fmt.Printf("\tLastBlockHash: %s\n", bx.LastBlockHash)
  } // endif no previous hash.
  fmt.Printf("\tMerkleHash: %s\n\n", bx.MerkleHash)

} // end func PrintBlockHeader.



// Print a single Transaction of the above types.
func PrintTransaction(gx Transaction, sx string) {
  fmt.Printf("%s: %d: %s\n", sx, gx.Tid, gx.Ttyp)
  PrintCoins(gx.Inputs, "OwnerIn ")
  PrintCoins(gx.Outputs, "OwnerOut ")
  fmt.Println("Tsig[:20]: ", gx.Tsig[:20])
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

// Summarize a tx on one line.
func PrintSummaryTx(ix int, tran Transaction) {

  fmt.Printf("[%d] %013d/%02d %s %8d %8d %8d %8d %8s\n", ix, tran.Tid, ix, tran.Ttyp, len(tran.Inputs), CoinValues(tran.Inputs), len(tran.Outputs), CoinValues(tran.Outputs), tran.Outputs[0].Owner[:5])

} // end PrintSummaryTx.

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
} // end func PrintMyCoins.

// Print all coins in a coin transaction.
func PrintCoins(pcoins []Coin, InOrOut string) {
  fmt.Printf("[ ")
  for px := 0; px < len(pcoins); px++ {
    PrintCoin(pcoins[px], InOrOut)
  } // end for px.
  fmt.Printf("%s Total: %d\n", InOrOut, CoinValues(pcoins))
} // end func PrintCoins.

// Print the Id, denomination and owner of a single coin.
func PrintCoin(thecoin Coin, InOrOut string) {
  if strings.Contains(InOrOut, "OwnerIn") {
    inhash := himitsu.DERToHash(thecoin.Owner)
    fmt.Printf("   %s: Id/Seq: %013d/%s Denom: %s: %s\n", InOrOut, thecoin.Cid, thecoin.Seq, thecoin.Denom, inhash)
  } else {
    fmt.Printf("   %s: Id/Seq: %013d/%s Denom: %s: %s\n", InOrOut, thecoin.Cid, thecoin.Seq, thecoin.Denom, thecoin.Owner)
  } // endinf Change Der to Hash.
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

  if len(thismap) == 0 {
    fmt.Printf("%s: map empty.\n", thistitle)
    return
  } // endif len map.

  for key, value := range thismap {
    fmt.Println("\t", key, "=", value)
  } // for thismap.

} // end func PrintMap.

//Print the Values of the  given Map.
func PrintMapValues(thismap map[string]string) {

  if len(thismap) == 0 {
    fmt.Printf("Map Empty.\n")
    return
  } // endif len map.

  for _, value := range thismap {
    fmt.Println(value)
  } // for thismap.

} // end func PrintMapValues.

func IntMap(mintmap map[string]int, thistitle string) {
  fmt.Println(thistitle)

  if len(mintmap) == 0 {
    fmt.Printf("%s: map empty.\n", thistitle)
    return
  } // endif len map.

  for key, value := range mintmap {
    fmt.Printf(",%s=%d", key[:5], value)
  } // for thismap.
  fmt.Println("")

} // end func IntMap.

func MintMap(mintmap map[string]int, thistitle string) {
  fmt.Println(thistitle)

  if len(mintmap) == 0 {
    fmt.Printf("%s: map empty.\n", thistitle)
    return
  } // endif len map.

  for key, value := range mintmap {
    fmt.Printf("%s=%d\n", key, value)
  } // for thismap.
  fmt.Println("")

} // end func IntMap.

