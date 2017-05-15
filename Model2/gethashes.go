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
  April 12, 2017
*/

/* DOCSTRING
 Usage:
  ./gethashes paycoins.conf

Get last hash functions from previous blocks, if any.
*/


import (
  "currency/methods"
  "currency/structures"
  "flag"
  "fmt"
  "io/ioutil"
  "os"
  "sort"
  "strings"
)

// Only runs if ledger exists (else run genesis)
// and newtxs file contains PayCoins transactions.
func main() {

  // Get Configs.
  flag.Parse()
  flay := flag.Args()
  if len(flay) == 0 {
    fmt.Println("Usage: lasthash <config-file>")
    os.Exit(1)
  } // if flay.

  // Load config values:
  confile := fmt.Sprintf("conf/%s", flay[0])
  mconf := methods.GetConfigs(confile)

  lastash := GetLastHash(mconf)
  fmt.Printf("Last hash retrieved: %s\n", lastash)

} // end func main.


// If there is a previous hash, get it from file.
// otherwise, return empty string.
func GetLastHash(smap map[string]string) string {
  hashes := GetDir(smap["blocks"], smap["hashes"])
  structures.Stringslice(hashes, "BlockHash files:")

  sort.Strings(hashes)
  lasthashfile := hashes[len(hashes)-1]
  fmt.Printf("Hashfile is %s\n", lasthashfile)
  lasthash := methods.ReadHash(lasthashfile)
  return lasthash

} // end func GetLastHash.



// List the files in a given directory, with a given suffix.
func GetDir(dirname string, suffix string) []string {
  fullpath := ""
  fileslice := []string{}
  files, err := ioutil.ReadDir(dirname)
  if err != nil {
    fmt.Printf("No files in %s.\n", dirname)
    return fileslice
  } // endif.

  for afile := range files {
    if strings.Contains(files[afile].Name(), suffix) {
      fullpath = dirname + "/" + files[afile].Name()
      fmt.Printf("%s in %s\n", suffix, fullpath)
      fileslice = append(fileslice, fullpath)
    } // endif.
  } // end for.
  return fileslice
} // end getdir.

