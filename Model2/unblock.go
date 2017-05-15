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
  April 6, 2017
*/

/* DOCSTRING
Usage:
 ./unblock blocks <file-prefix>

For each file in the blocks directory containing a numbered block, place it
into a map, unpack the block header, and print.
*/

import (
  "flag"
  "fmt"
  "os"
  "currency/himitsu"
  "currency/methods"
  "currency/structures"
)


// Get the blocks directory and unpack blocks files
func main() {

  levels := false
  flag.Parse()
  flay := flag.Args()
  if len(flay) == 0 || len(flay) > 3 {
    fmt.Println("Usage: unblock <blocksdir> block [verbose]")
    os.Exit(1)
  } // end if flay.

  blocks := methods.GetDir(flay[0], flay[1])

  for fx := 0; fx < len(blocks); fx++ {
    bloke, mermap := UnpackBlock(blocks[fx])
    structures.PrintBlockHeader(bloke)
    if len(flay) > 2 {
      if flay[2] == "levels" { levels = true }
      MerkWalk(bloke.MerkleHash, mermap, levels)
    } // endif verbose.
  } // end for blocks files.

} // end main.



// Get the contents of the 'block' file and extract the Block
// header.
func UnpackBlock(fn string) (structures.Block, map[string]string) {
  merkplus := map[string]string{}
  blocker := methods.ReadLedger(fn)
  blockhash := ""

  for ix := 0; ix < len(blocker); ix++ {
    if ix == 0 {
      blockhash = blocker[ix]
      fmt.Printf("ThisBlockHash: %s\n", blockhash)
    } else {
      // hashstring := himitsu.Hashit(himitsu.Hashit(blocker[ix]))
      hashstring := himitsu.Hashit(blocker[ix])
      merkplus[hashstring] = blocker[ix]
    } // endif blockhash.
  } // end for blocker.

  return methods.UnpackBlock(merkplus[blockhash]), merkplus

} // end func UnpackBlock.


// Walk the Merkle Map, given the root hash.
func MerkWalk(mroot string, mmap map[string]string, lev bool) {
  // fmt.Printf("In MerkWalk with: %s\n", mroot)
  mrlen := 44 // hash length.

  // If it's not two concatenated hashes, it is a transaction.
  if len(mmap[mroot]) != 2*mrlen {
    mtrx := methods.UnpackTransact(mmap[mroot])
    structures.PrintTransaction(mtrx, "MerkWalk")
    fmt.Println("")
  } else {
    lefthash := mmap[mroot][:mrlen]
    righthash := mmap[mroot][mrlen:]
    if lev {
      fmt.Printf("Left branch: %s = %s\n", lefthash, mmap[lefthash])
      fmt.Printf("Right branch: %s = %s\n", righthash, mmap[righthash])
    } // endif print intermediate levels.
    MerkWalk(lefthash, mmap, lev)
    MerkWalk(righthash, mmap, lev)
  } // endif leaf.

} // end func MerkWalk.


