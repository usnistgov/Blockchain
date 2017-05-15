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
)

// 01/31/2017 'scrooge.go' renamed to currency.go.
// It is an instance of a currency creating central bank,
// but it also performs Universal Basic Income Distribution.
// 02/13/2017 currency.go split into components:
// genesis.go, createcoins.go, ubipay.go, paycooins.go
// This is genesis.go.
// Test with: Ttyp == "BigBang", Outputs[0] == "v=HoR"

func main() {

  // Get Configs.
  flag.Parse()
  flay := flag.Args()
  if len(flay) == 0 {
    fmt.Println("Usage: genesis <config-file>")
    os.Exit(1)
  } // if flay.

  // Load config values:
  confile := fmt.Sprintf("conf/%s", flay[0])
  mainconf := methods.GetConfigs(confile)
  fmt.Println("\tCurrency Instance", methods.MilliNow(), "\n")
  structures.PrintMap(mainconf, "Config Values:")

  Genesis(mainconf)

} // end func main.


// Genesis creates the ledger.
func Genesis(confmap map[string]string) {
  ledger := []structures.Transaction{}

  // Load blockchain ledger or create a new genesis block.
  // Inventory extant coinbase and owners.
  _, err := os.Stat(confmap["ledge"])
  if err == nil {
    ledger = methods.LoadLedger(confmap["ledge"], ledger)
    fmt.Println("Txs in Ledger after Load: ", len(ledger))
    if confmap["verbose"] == "true" {
      structures.PrintLedger(ledger)
    } // endif verbose.
    fmt.Printf("\nLedger exists, length %d\n", len(ledger))
  } else {
    nextx := newtxes.CreateGenesis(confmap)
    ledger = append(ledger, nextx)
    structures.PrintTransaction(nextx, "Genesis:")
  } // end loadOrCreateLedger.

  // Publish the Blockchain:
  methods.StoreLedger(confmap["ledge"], ledger, "Blockchain")
  fmt.Printf("%d transactions stored.\n", len(ledger))

} // end func Genesis.


