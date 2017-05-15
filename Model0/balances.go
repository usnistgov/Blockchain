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
    "fmt"
    "currency/methods"
    "currency/structures"
    "strconv"
)

func main() {

  allbals := Balances("blocks/ledger")
  total := 0

  for pubkey := range(allbals) {
    fmt.Printf("Balance for %s = %d.\n", pubkey[:5], allbals[pubkey])
    total += allbals[pubkey]
  } // end for balmap.
  fmt.Printf("Total CoinBase = %d\n", total)

} // end func main.


// Get all coins and their owners, and calculate
//the balance of every account holder.
func Balances(ledgefile string) map[string]int {
    
  ledger := []structures.Transaction{}
  m1 := []structures.Coin{}
  thebals := map[string]int{}

  ledger = methods.LoadLedger(ledgefile, ledger)
  m1 = methods.M1(ledger)
  for bx := 0; bx < len(m1); bx++ {
    if _, ok := thebals[m1[bx].Owner]; ok {
      value, _ := strconv.Atoi(m1[bx].Denom)
      thebals[m1[bx].Owner] += value
    } else {
      value, _ := strconv.Atoi(m1[bx].Denom)
      thebals[m1[bx].Owner] = value
    } //endif haspkey.
 
  } // end for m1.

  return thebals

} // end func Balances.


