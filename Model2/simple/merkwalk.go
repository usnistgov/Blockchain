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
  April 2, 2017
*/

/* DOCSTRING
 Usage:
  ./merkwalk  map.txt

MerkWalk reads the hashes and transactions from a file to reconstitute a
Merkle Map. The Transactions are printed out in sequence.
*/

import (
  "currency/himitsu"
  "currency/methods"
  "currency/structures"
  "flag"
  "fmt"
  "os"
)


// Read in and reconstitute the Merkle Tree.
// Print it out.
func main() {

  flag.Parse()
  flay := flag.Args()
  if len(flay) != 1 {
    fmt.Println("Usage: merkwalk blocks/merk<seq>")
    os.Exit(1)
  } // endif flay.

  flatmap := methods.ReadLedger(flay[0])
  merkroot := flatmap[0]
  fmt.Printf("Merkle Root Hash: %s (len[%d])\n\n", merkroot, len(merkroot))
  flatmap = flatmap[1:]

  MerkleMap := MerkMake(flatmap)
  MerkWalk(merkroot, MerkleMap)

} // end main.


// Reconstitute the Merkle Tree.
func MerkMake(mapfile []string) map[string]string {
  mermap := map[string]string{}

  for lx := 0; lx < len(mapfile); lx++ {
    hashthis := himitsu.Hashit(mapfile[lx])
    mermap[hashthis] = mapfile[lx]
  } // end for mapfile.

  return mermap

} // end func MerkMake.

// Walk the Merkle Map, given the root hash.
func MerkWalk(mroot string, mmap map[string]string) {
  // fmt.Printf("In MerkWalk with: %s\n", mroot)
  mrlen := 44 // hash length.

  // If it's not two concatenated hashes, it is a transaction.
  if len(mmap[mroot]) != 2*mrlen {
    mtrx := methods.UnpackTransact(mmap[mroot])
    structures.PrintTransaction(mtrx, "MerkWalk")
    fmt.Println("")
  } else {
    lefthash := mmap[mroot][:mrlen]
    // fmt.Printf("Left hash: %s\n", lefthash)
    righthash := mmap[mroot][mrlen:]
    // fmt.Printf("Right Hash: %s\n", righthash)
    MerkWalk(lefthash, mmap)
    MerkWalk(righthash, mmap)
  } // endif leaf.

} // end func MerkWalk.

