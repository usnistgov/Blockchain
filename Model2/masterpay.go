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
  March 15, 2017
*/

/* DOCSTRING
 Usage:
  ./masterpay bufferconf.conf

Masterpay is the transaction listener that receives requests from the coin
client for Balances or PayCoins, validates them and stores them as new
transactions (newtxs) for periodic incorporation in the ledger by the
paycoins function. Masterpay and Paycoins hence communicate through the ledger,
which is created and validated by paycoins, and the newtxs, which are created
and validated by Masterpay.
*/

import (
  "bufio"
  "encoding/base64"
  "encoding/json"
  "flag"
  "fmt"
  "log"
  "net"
  "os"
  "currency/himitsu"
  "currency/methods"
  "currency/structures"
  "strconv"
  "strings"
  "time"
)

// File Scope for Balances:
var utxos = []structures.Coin{}

// Created from simple/echoserver.go Jan 21, 2017.
// Import stuff from masterpay.go to create the
// full scrooge listener.
func main() {

    fmt.Println(time.Now(), "\nListener v2.0")
    fmt.Println("Warning: You must set a max buffer length, else masterpay fails.")
    flag.Parse()
    flay := flag.Args()

    if len(flay) == 0 {
        fmt.Println("Usage: masterpay <config> [&]")
        os.Exit(1)
    } // endif flay.

    confile := fmt.Sprintf("conf/%s", flay[0])
    mmap := methods.GetConfigs(confile)
    bufmax, _ := strconv.Atoi(mmap["bufmax"])
    structures.PrintMap(mmap, "Config Values:")
    ledger := []structures.Transaction{}
    utxos := []structures.Coin{}
    fmt.Printf("Load Ledger from %s\n", mmap["ledge"])
    ledger = methods.LoadLedger(mmap["ledge"], ledger)
    fmt.Printf("Load LastUtxos from %s\n", mmap["lastutxos"])
    utxos = methods.GetLastUtxos(mmap)
    utxos = methods.M1UnUpdates(utxos, ledger)

    if len(utxos) == 0 {
      fmt.Println("FAIL: No Coins in Utxos.")
      os.Exit(1)
    } else {
      structures.PrintCoins(utxos, "StartingUtxos")
    } // endif get Utxos.

    listener, err := net.Listen("tcp", ":8081")
    if nil != err {
        log.Fatalln(err)
    }
    for {
        conn, err := listener.Accept()
        if nil != err {
            log.Fatalln(err)
        }
        Handle(mmap, conn, bufmax)
    }
}


// One connection for one transaction.
// Unwrap the received transaction,
// verify it, process it and
// send back the reply.
// Then close the connection.
func Handle(config map[string]string, conn net.Conn, maxbuf int) {
    defer conn.Close()
    p := make([]byte, maxbuf)
    processed := ""
    for {
        n, err := conn.Read(p)
        if err != nil {
            // fmt.Println(conn.RemoteAddr(), err)
            break
        }
        tx := UnwrapMessage(string(p[:n]))
        if n == maxbuf {
          processed = BuildErrorReply(config, tx)
        } else {
          processed = ParseTransaction(config, tx)
        } // endif ParseTrans.

        if strings.Contains(processed, "Error") {
          processed = BuildErrorReply(config, tx)
        } // endif error.
        if _, err := conn.Write([]byte(processed + "\n")); nil != err {
            // log.Println(conn.RemoteAddr(), err)
            break
        }
    }
}

// UnwrapMessage: base64 decode and Unmarshal.
func UnwrapMessage(msg string) structures.Transaction {
  bx := structures.Transaction{}
  bmsg, _ := base64.StdEncoding.DecodeString(msg)
  err := json.Unmarshal([]byte(bmsg), &bx)
  fmt.Printf("Masterpay: Message length=%d, unpacked length=%d\n", len(msg), len(bmsg))
  methods.CheckErrorInst(0, err, true)

  return bx

} // end func UnwrapMessage.

// ParseTransaction: distinguish Quit, Balance and PayCoins.
func ParseTransaction(conf map[string]string, px structures.Transaction) string {

  parsee := "undefined"
  rerr := methods.VerifyTransaction(px, false)
  methods.CheckErrorInst(1, rerr, true)
  if rerr != nil { fmt.Println("Verified? ", rerr) }

  switch(px.Ttyp) {
    case "Quit":
      scroogeair := himitsu.BaseDER(conf["pubkey"])
      pkin := px.Inputs[0]
      if scroogeair == pkin.Owner {
        fmt.Println("Masterpay Quitting: Authorized.")
      } else {
        fmt.Println("Masterpay Quitting: unauthorized.")
      } // endif scroogekey.
      os.Exit(1)

    case "Balance":
        if rerr == nil {
          firstop := px.Inputs[0]
          hashpub := himitsu.DERToHash(firstop.Owner)
          parsee = GiveBalance(conf, hashpub)
        } else {
            parsee = "Error"
        } // endif err.

    case "Transactions":
        if rerr == nil {
          firstop := px.Inputs[0]
          parsee = GiveTransactions(conf, firstop.Owner)
        } else {
            parsee = "Error"
        } // endif err.

    case "PayCoins":
      if rerr == nil {
        if msg, ok := methods.LegitimateTransaction(utxos, px); ok {
          utxos = methods.M1Update(utxos, px)
          fmt.Printf("Masterpay: %s\n", msg)
          parsee = Confirm("PayeeReceipt", conf, px, false)
        } else {
          fmt.Printf("Masterpay: %d %s\n", px.Tid, msg)
          structures.PrintTransaction(px, "ReferToDrawer")
          parsee = "ReferToDrawer"
        }
      } else {
        fmt.Println("Masterpay: PayCoins NOT Verified.")
        structures.PrintTransaction(px, "TxNotVerified")
        parsee = "TxVerificationError"
      } // endif error.

    case "Utxos":
      // Signal from Banker to update Utxos.
      // utxos = methods.LoadUtxos(conf["utxos"], utxos)

      ledger := []structures.Transaction{}
      ledger = methods.LoadLedger(conf["ledge"], ledger)
      utxos := methods.M1(ledger)
      structures.PrintCoins(utxos, "COINUPDATE")
      parsee = Confirm("Ack", conf, px, true)

    default:
      parsee = "Error"

  } // end switch Ttyp.

  return parsee

} // end func Parse Transaction.


// Filter my coins out of the ledger.
// mypub is the hash of my public key.
func MyCoins(fig map[string]string, mypub string) []structures.Coin {

    ledger := []structures.Transaction{}
    ledger = methods.LoadLedger(fig["ledge"], ledger)
    m1 := methods.M1(ledger)
    my1 := methods.GetMyCoins(m1, mypub)
    return my1

} // end func MyCoins.

// mypub is the DER form public key. But we need to search for the hash too.
func MyTransactions(fig map[string]string, mypub string) []structures.Coin {

    myledger := []structures.Transaction{}
    linearized := []structures.Coin{}
    myledger = methods.FilterLedger(fig["ledge"], myledger, mypub)
    linearized = LinearizeLedger(myledger)
    return linearized

} // end func MyTransactions.


// The coin client expects replies in the form of a coin structure.
// Each transaction should be linearized to put Transaction fields into coin form.
func LinearizeLedger(mytr []structures.Transaction) []structures.Coin {

  linear := []structures.Coin{}

  for mx := 0; mx < len(mytr); mx++ {

    lintr := mytr[mx]
    tidcoin := structures.Coin{ lintr.Tid, "0", "0", "" }
    linear = append (linear, tidcoin)
    ttypcoin := structures.Coin{ 0, "0", "0", lintr.Ttyp }
    linear = append (linear, ttypcoin)
    incoin := structures.Coin{ 0, "0", "0", "INPUTS:" }
    linear = append (linear, incoin)
    for cx := 0; cx < len(lintr.Inputs); cx++ {
      linear = append(linear, lintr.Inputs[cx])
    } // end for Inputs.
    outcoin := structures.Coin{ 0, "0", "0", "OUTPUTS:" }
    linear = append (linear, outcoin)
    for cx := 0; cx < len(lintr.Outputs); cx++ {
      linear = append(linear, lintr.Outputs[cx])
    } // end for Outputs.
    
  } // end for mytr.

  return linear

} // end func LinearizeLedger.


// GiveBalance: get the extant coin count for the given (hashed)
// public key and return it to the inquirer.
// Return wrapped as a signed (by scrooge) Transaction.
func GiveBalance(bconf map[string]string, pubki string) string {
  ledge := []structures.Transaction{}
  ledge = methods.LoadLedger(bconf["ledge"], ledge)
  utxos = methods.M1(ledge)
  outputs := methods.GetMyCoins(utxos, pubki)
  structures.PrintCoins(outputs, "MyBalance")
  ownair := himitsu.BaseDER(bconf["pubkey"])
  thinair := structures.Coin{0, "0", "0", ownair}
  inputs := []structures.Coin{ thinair }
  sig := methods.DoSign(bconf, outputs)
  tx := structures.Transaction{0, "Balance", inputs, outputs, sig}
  jsine, _ := json.Marshal(tx)
  ans := base64.StdEncoding.EncodeToString(jsine)
  return ans
} // end func GiveBalance.

// GiveTransactions: get the list of transactions for the given
// public key and return it to the inquirer.
// Return wrapped as a signed (by scrooge) Transaction.
func GiveTransactions(bconf map[string]string, pubki string) string {
  outputs := MyTransactions(bconf, pubki)
  fmt.Printf("%s has %d transactions.\n", pubki, len(outputs))
  ownair := himitsu.BaseDER(bconf["pubkey"])
  thinair := structures.Coin{0, "0", "0", ownair}
  inputs := []structures.Coin{ thinair }
  sig := methods.DoSign(bconf, outputs)
  tx := structures.Transaction{0, "Transactions", inputs, outputs, sig}
  jsine, _ := json.Marshal(tx)
  ans := base64.StdEncoding.EncodeToString(jsine)
  return ans
} // end func GiveTransactions.


func Confirm(receipt string, bconf map[string]string, txin structures.Transaction, ack bool) string {

  if txin.Ttyp != "Utxos" {
    wrapt := WrapTransaction(txin)
    SaveTransaction(bconf["newtxs"], wrapt, ack)
  } // endif anything but Utxos.

  ownair := himitsu.BaseDER(bconf["pubkey"])
  thinair := structures.Coin{0, "0", "0", ownair }
  inputs := []structures.Coin{ thinair }
  outputs := txin.Outputs
  sig := methods.DoSign(bconf, outputs)
  tx := structures.Transaction{0, receipt, inputs, outputs, sig}
  if ack {
    structures.PrintTransaction(tx, receipt)
  } // endif ack.
  wrout := WrapTransaction(tx)
  return wrout

} // end func Confirm.

// Build a Reply with "Error" Ttype and sign as Scrooge.
func BuildErrorReply(econf map[string]string, txin structures.Transaction) string {

  ownair := himitsu.BaseDER(econf["pubkey"])
  thinair := structures.Coin{0, "0", "0", ownair }
  inputs := []structures.Coin{ thinair }
  outputs := txin.Outputs
  sig := methods.DoSign(econf, outputs)
  tx := structures.Transaction{0, "Error", inputs, outputs, sig}
  wrapt := WrapTransaction(tx)
  return wrapt

} // end func BuildErrorReply.


// WrapTransaction: Marshal and Base64Encode.
func WrapTransaction(txin structures.Transaction) string {
  jtx, err := json.Marshal(txin)
  methods.CheckErrorInst(2, err, true)
  bjtx := base64.StdEncoding.EncodeToString(jtx)
  return bjtx

} // end func WrapTransaction.

// Append this transaction to the newtxs file.
func SaveTransaction(txfile string, txout string, ackn bool) {

  if ackn {
    fmt.Println("Write to", txfile)
    fmt.Println(txout)
  } // endif ackn.

  fn, erro := os.OpenFile(txfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
  methods.CheckErrorInst(3, erro, true)
  writer := bufio.NewWriter(fn)
  defer fn.Close()
  fmt.Fprintln(writer, txout)
  writer.Flush()

} // end func SaveTransaction.

