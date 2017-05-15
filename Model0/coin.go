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
  "bufio"
  "encoding/base64"
  "encoding/json"
  "flag"
  "fmt"
  "math/rand"
  "net"
  "os"
  "currency/himitsu"
  "currency/methods"
  "currency/newtxes"
  "currency/structures"
  "strconv"
  "strings"
  "time"
)


// This is the client for the Scrooge masterpay system.
// Commands include:
// - Balance (get my balance)
// - PayCoins (Pay from my balance to Payee)
// - Quit (Remotely close the masterpay listener. Only Scrooge can do this).
// Usage:
// coin Balance conf/<user.conf>
// coin PayCoins conf/<user.conf> <CoinCount> keys/<payee.pub>
// coin Quit conf/scrooge.conf

// This is the package scope for all balances received.
var ballons = map[string]int{}
var baltime = map[string]int64{}
var allaccounts = map[string][]structures.Coin{}
var rollover = 60000

// Package scope for pre-populating a random ramp.
var randramp = []int{}
var nexrand = 0

var row, col = 0, 0
var pairs = [][]string{
{"alice/bob", "bob/chaz", "chaz/dave", "dave/ellie", "ellie/fiona", "fiona/alice"}, 
{"alice/chaz", "bob/dave", "chaz/ellie", "dave/fiona", "ellie/alice", "fiona/bob"}, 
{"alice/dave", "bob/ellie", "chaz/fiona", "dave/alice", "ellie/bob", "fiona/chaz"}, 
{"alice/ellie", "bob/fiona", "chaz/alice", "dave/bob", "ellie/chaz", "fiona/dave"}, 
{"alice/fiona", "bob/alice", "chaz/bob", "dave/chaz", "ellie/dave", "fiona/ellie"}, 
 } // endpairs.

func main() {

  skey := "QSoT0u9VlnrL4wq2pppy+jB4lEbYJ7xeWnE1VKzgVic="
  correspondent := "127.0.0.1"
  flag.Parse()
  flay := flag.Args()
  // Seed the random ramp.
  randramp = Seedramp(100)

  if len(flay) < 2 {
    fmt.Println("Usage: coin Command Configfile")
    os.Exit(1)
  } // endif flags.

  cmap := ParseFlags(flay, skey)
  if cmap["cmd"] == "Error" {
    fmt.Printf(" '%s', # %s # \n", cmap["cmd"], cmap["strexp"])
    os.Exit(1)
  } // endif Error.

  // Individual commands, multiple transactions in a file or 
  // automatically generated transactions.
  switch cmap["cmd"] {

    case "Multitest":
      DoMultitest(cmap, skey, correspondent)

    case "Autotest":
      for ix:= 0; true; ix++ {
        DoAutotest(cmap, skey, correspondent)
        fmt.Printf("[%d] Wait till Clearing is completed ... ", ix)
        time.Sleep(time.Second*70)
        fmt.Println("  ... Start next round.")
      } // end forever.

    default:
      coins, expl := SendAMessage(cmap, correspondent)
      fmt.Println(expl)
      fmt.Printf("Output coins received: %d\n", structures.CoinValues(coins))

  } // end switch cmd.

} // end func main.


// Pre-populate the random ramp.
func Seedramp(orange int) []int {
  aramp := []int{}
  rand.Seed(int64(time.Now().UnixNano()))
  for ix := 0; ix < 1000; ix++ {
    nextrand := rand.Intn(orange)
    if nextrand == 0 { nextrand = 1 }
    aramp = append(aramp, nextrand) 
  } // end for 1000.
  return aramp
} // end func Seedramp.



// Bundle Multitest commands from main in here.
func DoMultitest(dmap map[string]string, fedkey string, corresp string) {

  cmdarray := ReadTests(dmap["scriptfile"])
  coins := []structures.Coin{}
  expl := ""

    for ix := 0; ix <len(cmdarray); ix++ {
      fmt.Printf("[%d] %s\n", ix, cmdarray[ix])
      argy := strings.Split(string(cmdarray[ix]), string(" "))
      mmap := ParseFlags(argy, fedkey)
      mmap["originalcmd"] = dmap["cmd"]
      coins, expl = SendAMessage(mmap, corresp)
      fmt.Println(expl)
      structures.PrintCoins(coins, "OwnerOut")
      time.Sleep(time.Millisecond)
    } // end for cmdarray.

} // end func DoMultitest.


// Generate an array of Autotest commands.
func DoAutotest(dmap map[string]string, fedkey string, corresp string) {

  // Autoconstruct builds the command list.
  paycmds := Autoconstructor(randramp, corresp)

  // execute the command list.
  for ix := 0; ix < len(paycmds); ix++ {
    fmt.Printf("[%d] %s\n", ix, paycmds[ix])
    argy := strings.Split(string(paycmds[ix]), string(" "))
    mmap := ParseFlags(argy, fedkey)
    mmap["originalcmd"] = dmap["cmd"]
    coins, expl := SendAMessage(mmap, corresp)
    fmt.Println(expl)
    fmt.Printf("Confirming %d coins sent.\n", structures.CoinCount(coins))
    time.Sleep(time.Second)
  }

} // end func DoAutotest.



func Autoconstructor(span []int, crosp string) []string {
  cmds := []string{}

  rand.Seed(int64(time.Now().Unix()))
  for ix := 0; ix < 6; ix++ {
    randone := Randy(span)
    sender, receiver := ReturnPair()
    randome := GetBalance(sender, crosp)
    payable := int(randome * randone / 100)
    if payable == 0 { payable = 1 }
    fmt.Printf("payable=%d, randome=%d, randone=%d, lenght=%d\n", payable, randome, randone, len(span))
    onecmd := fmt.Sprintf("PayCoins %s %d %s", sender, payable, receiver)
    cmds = append(cmds, onecmd)
  } // end for 120.

  return cmds

} // end func Autoconstructor.

func Randy(randels []int) int {
  return randels[rand.Intn(len(randels))]
} // end func Randy.

// ReturnPair: return a sender and receiver pair of names.
func ReturnPair() (string, string) {
  senrec := strings.Split(pairs[row][col], string("/"))
  sen := fmt.Sprintf("users/%s.conf", senrec[0])
  rec := fmt.Sprintf("users/%s.pub", senrec[1])
  col += 1
  if col == 6 { 
    col = 0; row += 1
    if row == 5 { row = 0 }
  }
  return sen, rec
} // end func ReturnPair.


// Return true if element is in slice, false otherwise.
func Contains(inslice []int, element int) bool {

  for ix := 0; ix < len(inslice); ix++ {
    if inslice[ix] == element {
      return true
    } // endif inslice.
  } // end for inslice.

  return false

} // end func Contains.


// Given a config file, build the Balance command, send it and get the reply.
func GetBalance(clifile string, crosp string) int {
  amap := methods.GetConfigs(clifile)
  amap["cmd"] = "Balance"
  amap["who"] = clifile
  coins, _ := SendAMessage(amap, crosp)
  myhash := himitsu.HashPublicKey(amap["pubkey"])
  ballons[myhash] = len(coins)
  baltime[myhash] = methods.MilliNow()
  allaccounts[myhash] = coins
  fmt.Printf("Balance for %s = %d.\n", clifile, len(coins))
  return len(coins)
} // end func GetBalance.




// SendAMessage: if command is PayCoins, do a GetBalance
// first. Only Pay coins if Balance exceeds payment.
// Return the balance and a string explanation.
func SendAMessage(smap map[string]string, corresp string) ([]structures.Coin, string) {

  mycoign := []structures.Coin{}
  rx := structures.Transaction{}
  reply := ""
  expm := ""
  pmap := make(map[string]string)
  for k, v := range smap {
    pmap[k] = v
  } // end copy map.

  if smap["cmd"] == "PayCoins" {

    // only do a Getbalance if not recently done.
    getnew := true; getmine := 0
    mash := himitsu.HashPublicKey(smap["pubkey"])

    if _, ok := ballons[mash]; ok {
      getnew = false
      timenow := methods.MilliNow()
      if timenow - baltime[mash] > int64(rollover) {
        getnew = true
      } // endif rollover.
    } // endif ballons.

    if getnew {
      pmap["cmd"] = "Balance"
      mycoign, expm = SendAMessage(pmap, corresp)
      getmine = structures.CoinValues(mycoign)
    } else {
      fmt.Println("Retrieve old balance for ", mash)
      getmine = ballons[mash]
      mycoign = allaccounts[mash]
    } // endif getbalance.

    // If there is enough, do PayCoins, else bail out.
    pam, _ := strconv.Atoi(smap["amount"])
    if getmine < pam {
      return mycoign, "Insufficient Funds: Payment Blocked." + expm
    } else {
      fmt.Printf("%s: My (%s) Balance = %d\n", expm, smap["who"], getmine)
    } // endif low balance.

  } // endif PayCoins exception.

  // The main body of SendAMessage starts here (after recursive exceptions).
  tosend := BuildMessage(smap, mycoign)
  if strings.Contains(tosend, "Error") {
    return mycoign, "Coin: Failed to Execute: " + smap["cmd"]
  } // endif error.

  // Send and receive message here:
  fmt.Println("Send now: ", smap["cmd"])
  reply = Transact(tosend, corresp, smap["cmd"])
  rx = UnwrapResult(reply)
  rerr := methods.VerifyTransaction(rx, false)
  methods.CheckErrorInst(0, rerr, true)
  fmt.Println("Verified? ", rerr)
  reply = "Transaction Completed."
  return rx.Outputs, reply

} // end func SendAMessage.


// ParseFlags: validate command and get configs.
func ParseFlags(flags []string, privilly string) map[string]string {

  cmap := make(map[string]string)
  cmd := "Error"; strexp := "OK."

  if CmdArrayHas(flags[0]) {
    cmd = flags[0]
  } else {
    cmd = "Error"
    strexp = "No such command: " + flags[0]
  } // endif cmd.

  if len(flags) > 1 {
    cmap = methods.GetConfigs(flags[1])
  } else {
    cmd = "Error"
    strexp = "Usage: ..."
  } // endif flags.

  switch(cmd) {

    case "Quit":
      if !AuthorizedBy(cmap["pubkey"], privilly) {
        cmd = "Error"
        strexp = "Unauthorized Quit instance."
      } // endif authorized.

    case "Balance":
      if methods.NoFile(flags[1]) {
        cmd = "Error"
        strexp = "Bad filename: " + flags[1]
      } // endif NoFile.

    case "Transactions":
      if methods.NoFile(flags[1]) {
        cmd = "Error"
        strexp = "Bad filename: " + flags[1]
      } // endif NoFile.

    case "Multitest":
      if methods.NoFile(flags[1]) {
        cmd = "Error"
        strexp = "Bad filename: " + flags[1]
      } else {
        cmap["scriptfile"] = flags[1]
      } // endif NoFile.

    case "Autotest":
      if methods.NoFile(flags[1]) {
        cmd = "Error"
        strexp = "Bad Directory Name: " + flags[1]
      } else {
        cmap["allpay"] = flags[1]
      } // endif NoFile.

    case "PayCoins":
      cmap["amount"] = flags[2]
      if methods.NoFile(flags[3]) {
        cmd = "Error"
        strexp = "Bad pubkey filename: " + flags[3]
      } else {
        cmap["payee"] = himitsu.HashPublicKey(flags[3])
      } // endif NoFile.

    case "CreateCoins":
      if !AuthorizedBy(cmap["pubkey"], privilly) {
        cmd = "Error"
        strexp = "Unauthorized Coin Creation."
      } // endif unauthorized.

      if methods.NoFile(flags[1]) {
        cmd = "Error"
        strexp = "Bad filename: " + flags[1]
      } // endif NoFile.

      // Coin specs are already in cmap["coin"] and cmap["denom"] so we're good.

    default:
      strexp = "Unknown command."
      cmd = "Error"

  } // end switch.

  cmap["cmd"] = cmd
  cmap["strexp"] = strexp 
  return cmap

} // end func ParseFlags.

// Open a text file, read it, put results in an array.
func ReadTests(filein string) []string {
  slices := []string{}
  if file, err := os.Open(filein); err != nil {
    fmt.Println("Error:", err)
    os.Exit(0)
  } else {
    defer file.Close()
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
      slices = append(slices, scanner.Text())
    } // end for scanner.
  } // endif file.

  return slices

} // end func ReadTests.


// BuildMessage: message is a Base64 string of a Balance, PayCoins or Quit
// command, constructed as a signed Transaction.
func BuildMessage(fmap map[string]string, coy []structures.Coin) string {

  switch(fmap["cmd"]) {
    case "Quit":
      ans, msg := BuildQuit(fmap)
      return CheckUm(ans, msg)

    case "Balance":
      ans, msg := BuildBalance(fmap)
      return CheckUm(ans, msg)

    case "Transactions":
      ans, msg := BuildClientCommand(fmap, fmap["cmd"])
      return CheckUm(ans, msg)

    case "PayCoins":
      ans, msg := BuildPay(fmap, coy)
      return CheckUm(ans, msg)

    case "CreateCoins":
      ans, msg := BuildCreate(fmap, coy)
      return CheckUm(ans, msg)

    default:
      return "Error: Unrecognized Command."

  } // end switch cmd.

} // end func BuildMessage.


//BuildQuit: Put the Quit command into a Transaction.
func BuildQuit(qmap map[string]string) (string, string) {

  bx := newtxes.GetQuit(qmap)
  fmt.Println("Sending Quit for", qmap["who"])
  jbx,_ := json.Marshal(bx)
  btrx := base64.StdEncoding.EncodeToString(jbx)
  return btrx, "Good:"

} // end func BuildQuit.

//BuildBalance: Put the Balance command into a Transaction.
func BuildBalance(qmap map[string]string) (string, string) {

  bx := newtxes.GetBalance(qmap)
  fmt.Println("Sending GetBalance for", qmap["who"])
  jbx, _ := json.Marshal(bx)
  btrx := base64.StdEncoding.EncodeToString(jbx)
  return btrx, "Good:"

} // end func BuildBalance.

// BuildClientCommand:  put the request for Commandname into a str.Transaction
func BuildClientCommand(qmap map[string]string, commandname string) (string, string) {

  bx := newtxes.GetClientCommand(qmap, commandname)
  fmt.Printf("Sending Get%s for %s\n", commandname, qmap["who"])
  jbx, _ := json.Marshal(bx)
  btrx := base64.StdEncoding.EncodeToString(jbx)
  return btrx, "Good:"

} // end func BuildTransactions.

//BuildPay: Put only the coins to send into the PayCoins Transaction.
func BuildPay(qmap map[string]string, allbal []structures.Coin) (string, string) {

  sendonly := []structures.Coin{}
  limit, _ := strconv.Atoi(qmap["amount"])
  fullpub := himitsu.BaseDER(qmap["pubkey"])
  for qx := 0; qx < limit; qx++ {
    nexcoin := allbal[qx]
    nexcoin.Owner = fullpub   // the hash is no good as input coin field.
    sendonly = append(sendonly, nexcoin)
  } // end for limit.

  fmt.Printf("BuildPay: %s pays %d coins:\n", qmap["who"], len(sendonly))
  trx := newtxes.PayCoins(qmap, sendonly, qmap["payee"])
  jtrx, _ := json.Marshal(trx)
  btrx := base64.StdEncoding.EncodeToString(jtrx)
  return btrx, "Good:"

} // end func BuildPay.

// BuildCreate: only scrooge can build and send the CreateCoins transaction.
func BuildCreate(qmap map[string]string, nobal []structures.Coin) (string, string) {

  newcoy := newtxes.CreateCoins(qmap)
  coincount, _ := strconv.Atoi(qmap["coin"])
  fmt.Printf("Creating %d New Coins.\n\n", coincount)
  ctrx, _ := json.Marshal(newcoy)
  bctrx := base64.StdEncoding.EncodeToString(ctrx)
  return bctrx, "Good:"

} // end func BuildCreate.



//CheckUm: if msg contains Error, return it, else return ans.
func CheckUm(good string, bad string) string {
  if strings.Contains(bad, "Error:") {
    return bad
  } else {
    return good
  } // endif bad.

} // end func CheckUm.



// Implement 'element in array' method.
func CmdArrayHas(candidate string) bool {
  CmdArray := []string{ "Balance", "Transactions", "Quit", "PayCoins", "CreateCoins", "Multitest", "Autotest" }
    for _, element := range CmdArray {
        if element == candidate {
            return true
        }
    }
    return false

} // end func CmdArrayHas.


// The hash of the pubkey must match the privileged string.
func AuthorizedBy(pubfile string, privileged string) bool {

  inkey := himitsu.HashPublicKey(pubfile)
  return inkey == privileged

} // end func AuthorizedBy.



// Transact: handle the tcp connection setup, message send
// reply handle and close.
func Transact(msgin string, partner string, comd string) string {

  // Get connection handle.
  conn, err := net.Dial("tcp", partner + ":8081")
  methods.CheckErrorInst(1, err, true)

  // send to socket
  fmt.Fprintf(conn, msgin + "\n")

  // Quit: nothing coming back, pull the plug.
  if comd == "Quit" {
    fmt.Println("\nListener going away.")
    os.Exit(1)
  } // endif Quit.

  // listen for reply
  conn.SetReadDeadline(time.Now().Add(time.Second*10))
  message, errr := bufio.NewReader(conn).ReadString('\n')
  methods.CheckErrorInst(2, errr, true)
   
  // Process reply.
  // fmt.Print("Listener reply: " + message)
  time.Sleep(time.Second)

  // Close the connection and  go away.
  conn.Close()

  return  message

} // end func Transact.


//UnwrapResult: Decode, Unmarshal and Print the reply.
func UnwrapResult(rez string) structures.Transaction {

  bx := structures.Transaction{}
  bsmsg, errb := base64.StdEncoding.DecodeString(rez)
  methods.CheckErrorInst(3, errb, true)
  err:= json.Unmarshal([]byte(bsmsg), &bx)
  methods.CheckErrorInst(4, err, true)
  // structures.PrintTransaction(bx, "\nReply Recvd:")
  return bx

} // end func UnwrapResult.



