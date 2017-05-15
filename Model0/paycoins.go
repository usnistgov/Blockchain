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
  "currency/methods"
  "currency/structures"
  "flag"
  "fmt"
  "os"
  "time"
)

// 01/31/2017 'scrooge.go' renamed to currency.go.
// It is an instance of a currency creating central bank,
// but it also performs Universal Basic Income Distribution.
// currency is split into genesis, createcoins, ubi, paycoins.
// This is paycoins.
// Only runs if ledger exists (else run genesis)
// and newtxs file contains PayCoins transactions.
func main() {

  // Get Configs.
  flag.Parse()
  flay := flag.Args()
  if len(flay) == 0 {
    fmt.Println("Usage: paycoins <config-file>")
    os.Exit(1)
  } // if flay.

  // Load config values:
  confile := fmt.Sprintf("conf/%s", flay[0])
  mainconf := methods.GetConfigs(confile)
  fmt.Println("\tScrooge Coin Instance", methods.MilliNow(), "\n")
  structures.PrintMap(mainconf, "Config Values:")
  timeable, err := time.ParseDuration(mainconf["timer"])
  methods.CheckErrorInst(1, err, true)

  if timeable == 0 {
    Paycoins(mainconf)
  } else {
    for { // ever
      Paycoins(mainconf)
      time.Sleep(timeable)
    } // end forever.
  } // endif timeable.

} // end func main.


// Paycoins validates new transactions from the client(s)
// and adds them to the ledger.
func Paycoins(confmap map[string]string) {
  ledger := []structures.Transaction{}
  m1 := []structures.Coin{}

  // Load blockchain ledger  if exists.
  _, err := os.Stat(confmap["ledge"])
  if err == nil {
    fmt.Println("\nLoad Blockchain:")
    ledger = methods.LoadLedger(confmap["ledge"], ledger)
    fmt.Println("Txs in Ledger after Load: ", len(ledger))
    if confmap["verbose"] == "true" {
      structures.PrintLedger(ledger)
      structures.PrintCoins(m1, "OwnerOut")
    } // endif verbose.
  } else {
    fmt.Println("Paycoins: no Ledger. Create one with genesis and add coins with createcoins")
    os.Exit(1)
  } // end loadLedger.

  // Process paycoins transactions.
  goodones, badones := methods.ProcessNewTransactions(confmap, ledger)
  if len(goodones) > 0 {
    for gx := 0; gx < len(goodones); gx++ {
      structures.PrintTransaction(goodones[gx], "NEW IN:")
      ledger = append(ledger, goodones[gx])
      fmt.Printf("[%d] %d Txs in Ledger after PNT, %d in rejects.\n", gx, len(ledger), len(badones))
    } // end for goodones.
  } // endif goodones.

  // Publish the Blockchain:
  methods.StoreLedger(confmap["ledge"], ledger, "Blockchain")
  methods.StoreLedger(confmap["rejects"], badones, "Rejects")
  fmt.Printf("%d transactions stored.\n", len(ledger))
  fmt.Printf("%d transactions rejected.\n", len(badones))
  if len(goodones) > 0 || len(badones) > 0 {
    err = os.Remove(confmap["newtxs"])
    methods.CheckErrorInst(0, err, true)
    fmt.Println("newtxs file deleted.")
  } // endif good'n'bad.

} // end func Currency.


