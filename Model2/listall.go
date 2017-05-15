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

/*
  Stephen Nightingale
  night@nist.gov
  NIST, Information Technology Laboratory
  March 1, 2017
*/

/* DOCSTRING
 Usage:
  ./listall blocks/ledger
  ./listall blocks/ledger [height]

Listall provides a full printout of every transaction in the ledger.
The height option limits the print to that point in the ledger.
*/

import (
  "currency/methods"
  "currency/structures"
  "flag"
  "fmt"
  "os"
  "strconv"
)


// Print all Transactions in the ledger or newtxes.
func main() {

  flag.Parse()
  flay := flag.Args()
  if len(flay) < 1 || len(flay) > 2 {
    fmt.Println("Usage: listall ledger/newtxes [height]")
    os.Exit(1)
  } // endif flay.

  ledger := []structures.Transaction{}
  ledger = methods.LoadLedger(flay[0], ledger)
  if len(flay) == 2 {
    height, _ := strconv.Atoi(flay[1])
    ledger = ledger[:height]
  }
  fmt.Printf("\n%d Transactions in %s.\n", len(ledger), flay[0]) 
  for ix := 0; ix < len(ledger); ix++ {
    fmt.Printf("\n[%d] ", ix)
    structures.PrintTransaction(ledger[ix], "Ledger:")
  } // end for.
  fmt.Printf("\n%d Transactions in %s.\n", len(ledger), flay[0]) 

} // end main.


