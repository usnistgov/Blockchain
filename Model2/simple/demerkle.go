package main

import (
  "currency/himitsu"
  "currency/methods"
  "currency/structures"
  "flag"
  "fmt"
  "os"
)

func main() {

  flag.Parse()
  flay := flag.Args()
  if len(flay) != 1 {
    fmt.Println("Usage: demerkle <merkin>")
    os.Exit(1)
  } // endif flay.

  flatledge := methods.ReadLedger(flay[0])
  MerkleMap := map[string]string{}
  merkleroot := ""

  // The first line is the Merkle Root, Hash from here.
  for ix := 0; ix < len(flatledge); ix++ {
    if ix == 0 {
      merkleroot = flatledge[ix]
      continue
    } // end if root.

    hash := himitsu.Hashit(flatledge[ix])
    MerkleMap[hash] = flatledge[ix]

  } // end for makemap.

  fmt.Printf("MerkleRoot = %s\n", merkleroot)
  structures.PrintMap(MerkleMap, "MerkleMap")

}


