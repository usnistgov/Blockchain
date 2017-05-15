package main

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


import (
  "flag"
  "fmt"
  "os"
  "currency/methods"
  "currency/newtxes"
  "currency/structures"
  "strconv"
)

// 01/31/2017 'scrooge.go' renamed to currency.go.
// It is an instance of a currency creating central bank,
// but it also performs Universal Basic Income Distribution.
// 02/13/2017 currency.go split into genesis, createcoins, ubi, paycoins.
// This is createcoins.
// Create coins as stipulated in the config file.
// Create coins as transmitted in newtxes. Must be signed by banker.
func main() {

  // Get Configs.
  flag.Parse()
  flay := flag.Args()
  if len(flay) == 0 {
    fmt.Println("Usage: createcoins <config-file>")
    os.Exit(1)
  } // if flay.

  // Load config values:
  confile := fmt.Sprintf("conf/%s", flay[0])
  mainconf := methods.GetConfigs(confile)
  fmt.Println("\tCurrency Instance: Create Coins", methods.MilliNow(), "\n")
  structures.PrintMap(mainconf, "Config Values:")

  Createcoins(mainconf)

} // end func main.


// Createcoins generates coinage in the ledger.
func Createcoins(confmap map[string]string) {
  ledger := []structures.Transaction{}
  m1 := []structures.Coin{}

  // Fail if ledger does not exist. genesis.go creates it.
  _, err := os.Stat(confmap["ledge"])
  if err == nil {
    fmt.Println("\nLoad Blockchain:")
    ledger = methods.LoadLedger(confmap["ledge"], ledger)
    fmt.Println("Txs in Ledger after Load: ", len(ledger))
  } else {
    fmt.Println("Createcoins: No ledger. genesis.go creates it.")
    os.Exit(1)
  } // endif  err.

  // Create coins if in args.
  coincount, _ := strconv.Atoi(confmap["coin"])
  fmt.Printf("\nCreate %d Coins of denomination %s.\n", coincount, confmap["denom"])
  if coincount > 0 {
    nexxtx := newtxes.CreateCoins(confmap)
    ledger = append(ledger, nexxtx)
    structures.PrintTransaction(nexxtx, "Coins:")
    fmt.Println("Txs in Ledger after Create: ", len(ledger))
  } // end CreateCoins.

  m1 = methods.M1(ledger)
  fmt.Printf("\nCoinBase has %d coins.\n", len(m1))

  // Publish the Blockchain:
  methods.StoreLedger(confmap["ledge"], ledger, "Blockchain")
  fmt.Printf("%d transactions stored.\n", len(ledger))

} // end func Currency.


