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
// currency split into genesis, createcoins, ubi, paycoins.
// This is ubi. if there are enough coins in the coinbase,
// ubi distributes evenly amngst all known members.
func main() {

  // Get Configs.
  flag.Parse()
  flay := flag.Args()
  if len(flay) == 0 {
    fmt.Println("Usage: ubi <config-file>")
    os.Exit(1)
  } // if flay.

  // Load config values:
  confile := fmt.Sprintf("conf/%s", flay[0])
  mainconf := methods.GetConfigs(confile)
  fmt.Println("\tScrooge Coin Instance", methods.MilliNow(), "\n")
  structures.PrintMap(mainconf, "Config Values:")

  UBI(mainconf)

} // end func main.


// Load the ledger, check for outstanding coins, distribute among members.
func UBI(confmap map[string]string) {
  ledger := []structures.Transaction{}

  // Load blockchain ledger if it exists.
  _, err := os.Stat(confmap["ledge"])
  if err == nil {
    fmt.Println("\nLoad Ledger:")
    ledger = methods.LoadLedger(confmap["ledge"], ledger)
    fmt.Println("Txs in Ledger after Load: ", len(ledger))
  } else {
    fmt.Println("UBI: No ledger. genesis.go creates it.")
    os.Exit(1)
  } // endif  err.

  // Register UBI Recipients' public keys:
  fmt.Println("\nRegister UBI Recipients:")
  ubiers := methods.HashSigs(confmap["ubiers"], ".pub")
  if len(ubiers) == 0 {
    fmt.Println("Warning: No users to distribute Coinage to.")
    os.Exit(1)
  } // end no ubiers.

  for ubi := range ubiers {
    fmt.Println("\t", ubiers[ubi])
  } // end for ubi.

  // Check the _Unassigned_ Coinbase and Distribute UBI Payments:
  m0 := methods.M0(confmap, ledger)
  fmt.Printf("\nCoinBase has %d coins.\n", len(m0))
  structures.PrintCoins(m0, "OwnerOut")
  fmt.Println("\nDistribute UBI Payments.")
  ledgeplus := newtxes.DistributeUBI(confmap, m0, ubiers)
  structures.PrintLedger(ledgeplus)
  fmt.Printf("Adding %d new transactions to ledger (%d)\n", len(ledgeplus), len(ledger))
  for lx := 0; lx < len(ledgeplus); lx++ {
    ledger = append(ledger, ledgeplus[lx])
  } // end for ledgeplus.
  fmt.Println("Txs in Ledger after UBI: ", len(ledger))

  // Publish the Blockchain:
  methods.StoreLedger(confmap["ledge"], ledger, "Blockchain")
  fmt.Printf("%d transactions stored.\n", len(ledger))

} // end func UBI.


