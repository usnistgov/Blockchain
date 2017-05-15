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
  ./sum blocks/ledger

Sum summarizes every transaction in the ledger on one line, including the
number of inputs and their value, the number of outputs and their value,
the truncated hash of the output owner, and whether the transaction validates.
*/

import (
  "flag"
  "fmt"
  "os"
  "currency/methods"
  "currency/structures"
  // "time"
)


// Summarize Transactions in the ledger.
func main() {

  flag.Parse()
  flay := flag.Args()
  if len(flay) != 1 {
    fmt.Println("Usage: sum ledger/newtxes")
    os.Exit(1)
  } // endif flay.

  ledger := []structures.Transaction{}
  ledger = methods.LoadLedger(flay[0], ledger)
  PrintSummaryTxs(ledger)

} // end func main.

func PrintSummaryTxs(ledger []structures.Transaction) {

  fmt.Printf("Seq/TransactionId  Transact %8s %8s %8s %8s %8s %8s\n", "Ver",  "Inputs",  "Value",  "Outputs",  "Value",  "Recip")

  for ix := 0; ix < len(ledger); ix++ {
    very := "good."
    tran := ledger[ix]
    verf := methods.VerifyTransaction(tran, false)
    if verf == nil { very = "yes" } else { very = "no " }
    fmt.Printf("[%d] %013d/%02d %s %8s %8d %8d %8d %8d %8s\n", ix, tran.Tid, ix, tran.Ttyp, very, len(tran.Inputs), structures.CoinValues(tran.Inputs), len(tran.Outputs), structures.CoinValues(tran.Outputs), tran.Outputs[0].Owner[:5])
  } // end for.

  fmt.Printf("\n%d Transactions in the Ledger.\n",len(ledger)) 

} // end PrintSummaryTxs.


