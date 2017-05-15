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
  April 3, 2017
*/

/* DOCSTRING
 Usage:
  ./compress blocks/ledger <data.fl>

Compress reads in the ledger in Bas64 format, and writes it out compressed.
*/

import (
  "compress/flate"
  "currency/methods"
  "flag"
  "fmt"
  "os"
)

func main() {

  flag.Parse()
  flay := flag.Args()
  if len(flay) != 2 {
    fmt.Println("Usage: compress <ledgerin> <smallout>")
    os.Exit(1)
  } // endif flay.

  compress(flay[0], flay[1])
}


func compress(inputFile, outputFile string) {
  ledges := methods.ReadLedger(inputFile)
  o, _ := os.Create(outputFile)
  defer o.Close()
  f, _ := flate.NewWriter(o, flate.BestCompression)
  defer f.Close()

  for ix := 0; ix < len(ledges); ix++ {
    bytesout := []byte(fmt.Sprintf("%s\n", ledges[ix]))
    f.Write(bytesout)
  } // end for ledger.
}
