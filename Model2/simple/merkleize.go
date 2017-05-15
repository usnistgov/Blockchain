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
  March 31, 2017
*/

/* DOCSTRING
 Usage:
  ./merkleize blocks/ledger [length]

Merkleize creates a series of MerkleMaps from the given ledger, to the
given length, which must be an exponential of 2.
*/

import (
  "compress/flate"
  "currency/himitsu"
  "currency/methods"
  "flag"
  "fmt"
  "os"
  "strconv"
)


// Form the ledger into a Merkle Tree.
// Transactions beyond 2**N are left in a remainder list.
func main() {

  flag.Parse()
  flay := flag.Args()
  if len(flay) != 3 {
    fmt.Println("Usage: merkleize <length> ledger/newtxes <mapout.fl>")
    os.Exit(1)
  } // endif flay.

  exp, _ := strconv.Atoi(flay[0])
  ledger := methods.ReadLedger(flay[1])
  ledger = ledger[:exp]

  pyramids := [][]string{}
  hashlevel := []string{}
  MerkleMap := map[string]string{}
  levels := 0

  // Generate the data level hashes for the Merkle Map.
  for ix := 0; ix < len(ledger); ix++ {
    hash := himitsu.Hashit(ledger[ix])
    MerkleMap[hash] = ledger[ix]
    hashlevel = append(hashlevel, hash)
  } // end for.
  pyramids = append(pyramids, hashlevel)

  steps := []string{}
  for levels = 0; len(pyramids[levels]) > 1; levels++ {
    steps = []string{}
    for hx := 0; hx < len(pyramids[levels])-1; hx += 2 {
      concat := pyramids[levels][hx] + pyramids[levels][hx+1]
      hashhash := himitsu.Hashit(concat)
      MerkleMap[hashhash] = concat
      steps = append(steps, hashhash)
    } // end for hash level.
    pyramids = append(pyramids, steps)
  } // end for levels.

  // Print the MerkleMap and the Hash Levels ('pyramids')
  fmt.Println(pyramids[levels][0])
  compress(MerkleMap, flay[2])
  // structures.PrintMapValues(MerkleMap)


} // end main.


func compress(merkin map[string]string, outputFile string) {
  o, _ := os.Create(outputFile)
  defer o.Close()
  f, _ := flate.NewWriter(o, flate.BestCompression)
  defer f.Close()

  for _, value := range merkin {
    bytesout := []byte(fmt.Sprintf("%s\n", value))
    f.Write(bytesout)
  } // end for ledger.
}
