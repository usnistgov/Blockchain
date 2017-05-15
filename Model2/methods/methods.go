package methods

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
  March 13, 2017
*/

 
import (
  "bufio"
  "encoding/base64"
  "encoding/json"
  "fmt"
  "currency/himitsu"
  "currency/structures"
  "io/ioutil"
  "os"
  "sort"
  "strconv"
  "strings"
  "time"
)



// Check and Report Error Condition.
func CheckError(err error, bailif bool) {

      if err != nil {
        fmt.Println("Error: ", err)
        if bailif {
          os.Exit(1)
        } // endif balif.
      } // end if err.

} // end func CheckError.

func CheckErrorInst(instance int, err error, bailif bool) {

      if err != nil {
        fmt.Printf("[%d] Error: %s\n", instance, err)
        if bailif {
          os.Exit(1)
        } // endif balif.
      } // end if err.

} // end func CheckError.

// Return Unix time in Milliseconds.
func MilliNow() int64 {
  tim := time.Now()
  nanos := tim.UnixNano()
  millis := nanos/1000000
  return millis
} // end MilliNow.

// Read in the config file and store it in a map.
func GetConfigs(filename string) (mapvec map[string]string) {
  mapvec = make(map[string]string)
  argsin, _ := ioutil.ReadFile(filename)
  argvec := strings.Split(string(argsin), "\n")
  for line := range argvec {
    keyvl := strings.Split(argvec[line], "=")
    if len(keyvl) == 2 {
      mapvec[keyvl[0]] = keyvl[1]
    } // end len.
  } //end for line.

  // save the filename here, too.
  mapvec["who"] = filename 
  // And get the pubkey Hash and DER.
  if _, ok := mapvec["pubkey"]; ok {
    mapvec["pubhash"] = himitsu.HashPublicKey(mapvec["pubkey"])
    mapvec["pubder"] = himitsu.BaseDER(mapvec["pubkey"])
  } // endif pubkey present.
  if _, okay := mapvec["payee"]; okay {
    mapvec["payhash"] = himitsu.HashPublicKey(mapvec["payee"])
  } // endif payee present.
  return mapvec

} // getConfigs.

// List the files in a given directory, with a given suffix.
func GetDir(dirname string, suffix string) []string {
  var fullpath string
  fileslice := []string{}
  files, err := ioutil.ReadDir(dirname)
  if err != nil {
    fmt.Printf("No files in %s.\n", dirname)
    return fileslice
  } // endif.

  for afile := range files {
    if strings.Contains(files[afile].Name(), suffix) {
    fullpath = dirname + "/" + files[afile].Name()
      fileslice = append(fileslice, fullpath)
    } // endif.
  } // end for.
  return fileslice
} // end getdir.


// For all public key files in a directory,
// return a slice containing their hashes.
func HashSigs(dirname string, suffix string) (anonymice []string) {

  pubfiles := GetDir(dirname, suffix)
  for hfile := range pubfiles {
    sighash := himitsu.HashPublicKey(pubfiles[hfile])
    anonymice = append(anonymice, sighash)
  } // end for.
  return anonymice
} // end HashSigs.


// Load the ledger of Transactions from a given file.
// Return it as a list.
// TO BE MODIFIED TO: return a map of 'blocks' with hash pointers.
func LoadLedger(ledfile string, ledge []structures.Transaction) []structures.Transaction {

    file, err := os.Open(ledfile)
    CheckErrorInst(5, err, true)
    defer file.Close()

    scanner := bufio.NewScanner(file)

    for scanner.Scan() {             // internally, it advances token based on separator
        data, _ := base64.StdEncoding.DecodeString(scanner.Text())  // token in unicode-char
        var tx structures.Transaction
        err:= json.Unmarshal([]byte(data), &tx)
        CheckErrorInst(6, err, true)
        ledge = append(ledge, tx)

    } // end for scanner.

    return ledge

} // end func LoadLedger.




// Return the ledger in its Base64 (or byte string) form.
func ReadLedger(ledfile string) []string {

    file, err := os.Open(ledfile)
    CheckErrorInst(5, err, true)
    defer file.Close()

    scanner := bufio.NewScanner(file)
    ledge := []string{}

    for scanner.Scan() {
        ledge = append(ledge, scanner.Text())

    } // end for scanner.
    return ledge
} // end func ReadLedger.


// Return the value of a 1 line file.
func ReadHash(ledfile string) string {

    file, err := os.Open(ledfile)
    ledge := []string{}
    CheckErrorInst(6, err, true)
    defer file.Close()
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        ledge = append(ledge, scanner.Text())
    } // end for scanner.
    return ledge[len(ledge)-1]

} // end func ReadHash.



// Input is a Block Header in Base64. unpack and return it.
func UnpackBlock(trans string) structures.Block {

  data, _ := base64.StdEncoding.DecodeString(trans)
  var bx structures.Block
  err:= json.Unmarshal([]byte(data), &bx)
  CheckErrorInst(10, err, true)

  return bx

} // end func UnpackBlock.



// Input is a transaction in Base64. unpack and return it.
func UnpackTransact(trans string) structures.Transaction {

  data, _ := base64.StdEncoding.DecodeString(trans)
  var tx structures.Transaction
  err:= json.Unmarshal([]byte(data), &tx)
  CheckErrorInst(7, err, true)

  return tx

} // end func UnpackTransact.


// Input is a coin in Base64. Unpack and return it.
func UnpackCoin(cons string) structures.Coin {

  data, _ := base64.StdEncoding.DecodeString(cons)
  var ux structures.Coin
  err:= json.Unmarshal([]byte(data), &ux)
  CheckErrorInst(8, err, true)

  return ux

} // end func UnpackCoin.


// Update from the last saved utxos file.
func GetLastUtxos(smap map[string]string) []structures.Coin {
  utxfiles := GetDir(smap["blocks"], smap["lastutxos"])
  lutch := []structures.Coin{}

  if len(utxfiles) == 0 {
    return lutch
  } // endif no previous.

  sort.Strings(utxfiles)
  lastutxfile := utxfiles[len(utxfiles)-1]
  return LoadUtxos(lastutxfile, lutch)

} // end func GetLastUtxos.




// Load the Unspent coins (UTXOS) from the utxos file.
func LoadUtxos(utfile string, utch []structures.Coin) []structures.Coin {

  file, err := os.Open(utfile)
  CheckErrorInst(8, err, true)
  defer file.Close()
  scanner := bufio.NewScanner(file)

  for scanner.Scan() {
    data, _ := base64.StdEncoding.DecodeString(scanner.Text())
    var ux structures.Coin
    err := json.Unmarshal([]byte(data), &ux)
    CheckErrorInst(8, err, true)
    utch = append(utch, ux)
 } // end for scanner.

  return utch

} // end func LoadUtxos.


// Load the ledger of Transactions from a given file.
// Filter it for 'pubki' related transactions and return it as a list.
// TO BE MODIFIED TO: return a map of 'blocks' with hash pointers.
// Note: pubky is the DER form. We need the hash too.
func FilterLedger(ledfile string, myledge []structures.Transaction, pubky string) []structures.Transaction {

    file, err := os.Open(ledfile)
    CheckErrorInst(10, err, true)
    defer file.Close()
    hashky := himitsu.DERToHash(pubky)

    scanner := bufio.NewScanner(file)

    for scanner.Scan() {             // internally, it advances token based on separator
        data, _ := base64.StdEncoding.DecodeString(scanner.Text())  // token in unicode-char
        var tx structures.Transaction
        err:= json.Unmarshal([]byte(data), &tx)
        CheckErrorInst(11, err, true)
        if IAmInvolved(tx, pubky, hashky) {
          myledge = append(myledge, tx)
        } // endif IAminvolved.

    } // end for scanner.
    return myledge
} // end func FilTerLedger.


// Return true if Transaction Inputs contains pubkey or Outputs contains its hash.
func IAmInvolved(mytx structures.Transaction, derpub string, hashpub string) bool {
  innies := mytx.Inputs
  outies := mytx.Outputs

  for mx := 0; mx < len(innies); mx++ {
    if innies[mx].Owner == derpub {
      return true
    } // endif innies.
  } // end for innies.

  for mx := 0; mx < len(outies); mx++ {
    if outies[mx].Owner == hashpub {
      return true
    } // endif outies.
  } // end for outies.

  return false

} // end func IAmInvolved.





// Store the ledger back in ledger file.
func StoreLedger(ledfile string, lodge []structures.Transaction, herald string) {

  fp, _ := os.Create(ledfile)
  defer fp.Close()
  fmt.Printf("Store the %s.\n", herald)
  for ix := 0; ix < len(lodge); ix++ {
    jase, _ := json.Marshal(lodge[ix])
    dat := base64.StdEncoding.EncodeToString(jase)
    _, _ = fp.WriteString(dat + "\n")
    // bc, _ := fp.WriteString(dat + "\n")
    // fmt.Println(bc, " bytes  written.")
  } // end for print ledger.

} // end func StoreLedger.




// Store a block in the numbered block file.
func StoreBlock(blockhash string, ledfile string, thismap map[string]string) {

  if len(thismap) == 0 {
    fmt.Printf("Map Empty.\n")
    return
  } // endif len map.

  fp, _ := os.Create(ledfile)
  defer fp.Close()
  fmt.Printf("Store %s.\n", ledfile)
  _, _ = fp.WriteString(blockhash + "\n")

  for _, value := range thismap {
    _, _ = fp.WriteString(value + "\n")
  } // for thismap.

} // end func StoreBlock.



// Store a hash string as a single line file.
func StoreHash(ledfile string, blockhash string) {

  fp, _ := os.Create(ledfile)
  defer fp.Close()
  fmt.Printf("Store Hash in %s.\n", ledfile)
  _, _ = fp.WriteString(blockhash + "\n")

} // end func StoreHash.




// Store the Unspent Transactions.
func StoreUtxos(utfile string, utch []structures.Coin, herald string) {
  fp, _ := os.Create(utfile)
  defer fp.Close()
  fmt.Printf("Store the %s.\n", herald)
  for ix := 0; ix < len(utch); ix++ {
    jase, _ := json.Marshal(utch[ix])
    dat := base64.StdEncoding.EncodeToString(jase)
    _,_ = fp.WriteString(dat + "\n")
  } // end for utch.

} // end func StoreUtxos.



// Count Output Banker coins as 'M0'.
func M0(conf map[string]string, ledge []structures.Transaction) []structures.Coin {
  m1 := M1(ledge)
  m0 := []structures.Coin{}
  bankhash := himitsu.HashPublicKey(conf["pubkey"])

  for ix := 0; ix < len(m1); ix++ {
    if m1[ix].Owner == bankhash {
      m0 = append(m0, m1[ix])
    } // else ignore the rest.
  } // end for m1.

  return m0

} // end func M0,


// Enumerate the Coins within a CreateCoins transaction.
func CountCoins(onetrans []structures.Coin, allcoins []structures.Coin) []structures.Coin {

  for cx := 0; cx < len(onetrans); cx++ {
    allcoins = append(allcoins, onetrans[cx])
  } // end for cx.

  return allcoins

} // end CountCoins.




// Count all extant coins, and identify their current owners.
//Count input coins and output coins. If an outputs appears
// as a later input, remove it.
func M1(led []structures.Transaction) []structures.Coin {

  unspent := []structures.Coin{}

  for ix := 0; ix < len(led); ix++ {
    nx := led[ix]

    switch(nx.Ttyp) {

      case "BigBang":
        continue

      case "CreateCoins":
        unspent = append(unspent, nx.Outputs...)

      case "PayCoins":
        tryins := nx.Inputs
        tryouts := nx.Outputs
        for ty := 0; ty < len(tryins); ty++ {
          tx := Index(unspent, tryins[ty])
          if tx != -1 {
            if len(unspent) == 1 {
              unspent = []structures.Coin{} // zero out the current value.
            } else {
              unspent = append(unspent[:tx], unspent[tx+1:]...)
            } // endif len(1).
          } // endif tx // input coin is already spent.
        } // end for tryins.
        unspent = append(unspent, tryouts...)

      default:
        fmt.Printf("Bad transaction type: %s\n", nx.Ttyp)
        os.Exit(1)

    } // end switch Ttyp.

  } // end for led.

  return unspent

} // end func M1.


// Update Utxos from a block of transactions.
func M1UnUpdates(emone []structures.Coin, tranx []structures.Transaction) []structures.Coin {

  for ux := 0; ux < len(tranx); ux++ {

    thx := tranx[ux]
    if thx.Ttyp == "BigBang" { continue }
    if thx.Ttyp == "CreateCoins" {
      emone = append(emone, thx.Outputs...)
      continue
    } // endif CreateCoins.

    emone = M1Update(emone, thx)
  } // end for txs.

  return emone

} // end func M1UnUpdates.


// Update Utxos from a block of base64 wrapped transactions.
func M1Updates(emone []structures.Coin, txs []string) []structures.Coin {

  for ux := 0; ux < len(txs); ux++ {

    thx := UnpackTransact(txs[ux])
    if thx.Ttyp == "BigBang" { continue }
    if thx.Ttyp == "CreateCoins" {
      emone = append(emone, thx.Outputs...)
      continue
    } // endif CreateCoins.

    emone = M1Update(emone, thx)
  } // end for txs.

  return emone

} // end func M1updates.



// Update M1 unspent transaction outputs on a real-time basis.
func M1Update(emone []structures.Coin, intx structures.Transaction) []structures.Coin {

  newin := intx.Inputs
  newout := intx.Outputs
  emup := []structures.Coin{}
  flag := false

  for mx := 0; mx < len(emone); mx++ {
    if !TransContains(newin, emone[mx]) {
      emup = append(emup, emone[mx])
    } else {
      flag = true
    } // endif update emone.
  } // end for emone.

  if flag {
    emup = append(emup, newout...)
  } // endif add new outputs.

  return emup

} // end func M1Update.



// TransContains: if the utxo is in the current tx it's a new spend,
// so remove it.
func TransContains(inn []structures.Coin, utxo structures.Coin) bool {
  for nx := 0; nx < len(inn); nx++ {
    if utxo.Cid == inn[nx].Cid && utxo.Seq == inn[nx].Seq {
      return true
    } // endif
  } // end for inputcoins.

  return false

} // end func TransContains.




// Return the index of a coin within the slice.
func Index(coins []structures.Coin, acoin structures.Coin) int {
  for cx := 0; cx < len(coins); cx++ {
    if (coins[cx].Cid == acoin.Cid) && (coins[cx].Seq == acoin.Seq) {
      return cx
    } // endif samecoin.
  } // end for coins.
  return -1
} // end func Index.



// For every coin in every PayCoins transaction,
// find the Coin Id/Seq in originals, and
// substitute the new owner.
func ReassignedCoins(later []structures.Coin, earlier []structures.Coin) []structures.Coin {

  for ex := 0; ex < len(earlier); ex++ {
    for lx := 0; lx < len(later); lx++ {
      
      if later[lx].Cid == earlier[ex].Cid {
        if later[lx].Seq == earlier[ex].Seq {
          earlier[ex].Owner = later[lx].Owner
        } // endif Seq.
      } // endif Cid.

    } // end func later.
  } // end for earlier.

  return earlier

} // end func ReassignedCoins.



// Sign the Marshalled outputs and return the base64 encoded signature.
func DoSign(conmap map[string]string, outs []structures.Coin) string {
  signer, _ := himitsu.LoadPrivateKey(conmap["privkey"])
  jtx, _ := json.Marshal(outs)
  htx := himitsu.Bashbin(jtx)
  signtrn, _ := signer.Sign(htx)
  sig := base64.StdEncoding.EncodeToString(signtrn)
  return sig
} // end func DoSign.

// Loop through every Transaction, extract the owner's
// public key from the Inputs,the marshalled outputs
// and the signature. Verify. Or Not.
func VerifyTransactions(goodledge []structures.Transaction, noisy bool) bool {
  alltrue := true

  for ix := 0; ix < len(goodledge); ix++ {
    verified := VerifyTransaction(goodledge[ix], noisy)
    if verified == nil {
      fmt.Printf("%d verified.\n", goodledge[ix].Tid)
    } else {
      structures.PrintTransaction(goodledge[ix], goodledge[ix].Ttyp)
      fmt.Printf("%d failed.\n", goodledge[ix].Tid)
      alltrue = false
    } // end if verified.
  } // end for ix.

  return alltrue

} // end func VerifyTransactions.


// Verify one transaction.
func VerifyTransaction(onetran structures.Transaction, verbose bool) error {

  if verbose {
    structures.PrintTransaction(onetran, onetran.Ttyp)
  } // endif verbose.

  pubk := onetran.Inputs[0].Owner
  pubb, errb := base64.StdEncoding.DecodeString(pubk)
  CheckErrorInst(1, errb, true)
  verifier, errv := himitsu.ParseDERKey(pubb)
  CheckErrorInst(2, errv, true)
  signet, errs := json.Marshal(onetran.Outputs)
  CheckErrorInst(3, errs, true)
  signable := himitsu.Bashbin(signet)
  bsig, errbs := base64.StdEncoding.DecodeString(onetran.Tsig)
  CheckErrorInst(4, errbs, false)
  verified := verifier.Unsign(signable, bsig)

  return verified

} // end func VerifyTransaction.

// Get my coins.
func GetMyCoins(pcoins []structures.Coin, mypub string) []structures.Coin {
  justmine := []structures.Coin{}

  for px := 0; px < len(pcoins); px++ {
    if pcoins[px].Owner == mypub {
      justmine = append(justmine, pcoins[px])
    } // end if mypub.
  } // end for px.

  return justmine

} // end func PrintCoins.


// Process New Transactions, stored in 'newtxs' file by the masterpay listener.
// Check that transactions verify and are not double spends. Add good ones to
// ledger, send bad ones to rejects.
func ProcessNewTransactions(config map[string]string, lodger []structures.Transaction, emmy []structures.Coin) ([]structures.Transaction, []structures.Transaction, []structures.Coin) {

  newin := []structures.Transaction{}
  newout := []structures.Transaction{}
  rejects := []structures.Transaction{}

  if NoFile(config["newtxs"]) {
    fmt.Println("No New Transactions.")
    return newout, rejects, emmy
  } else {
    newin = LoadLedger(config["newtxs"], newin)
    fmt.Printf("%d New Transactions in.\n", len(newin))
  } // endif NoFile.

  // cycle through newtx sorting into good and bad.
  for ix := 0; ix < len(newin); ix++ {
    checking := newin[ix]
    structures.PrintTransaction(checking, "\nNEW TRANSACTION:")

    // does the signature verify.
    if VerifyTransaction(checking, false) != nil {
      rejects = append(rejects, checking)
      fmt.Printf("PNT: Tx [%d] failed to verify.\n", ix)
      continue
    } // endif Verify.

    // only expecting CreateCoins or PayCoins transactions.
    switch(checking.Ttyp) {

      case "CreateCoins":
        if BankerOnly(config, checking, false) {
          newout = append(newout, checking)
        } else {
          fmt.Printf("[%d] REJECTED: Only Banker can Create Coins.\n", ix)
          rejects = append(rejects, checking)
        } // endif BankerOnly.

      case "PayCoins":
        if msg, ok := LegitimateTransaction(emmy, checking); ok {
        // if NotDoubleSpent(emmy, checking) {
          fmt.Println(msg)
          newout = append(newout, checking)
          emmy = M1Update(emmy, checking)
        } else {
          fmt.Printf("[%d] REJECTED: Double Spend attempt.\n", ix)
          rejects = append(rejects, checking)
        } // endif NDS.

      default:
        fmt.Printf("[%d] REJECTED: Bad Transaction Type.\n", ix)
        rejects = append(rejects, checking)
    } // end switch Ttyp.

  } // end for newin.

  fmt.Printf("PNTs: %d Txs In, %d Good, %d Bad.\n", len(newin), len(newout), len(rejects))
  return newout, rejects, emmy

} // end func ProcessNewTransactions.



// LegitimateTransaction: weed out bad balance, spent txs, stolen coins.
// For every coin in the input transaction , it must be explicitly ruled in or out of the Utxos.
// So length returned from LegitTrans must be equal to the number of coins in the input Transaction.
func LegitimateTransaction(emmone []structures.Coin, sus structures.Transaction) (string, bool) {

  ins := sus.Inputs
  outs := sus.Outputs

  if !Balanced(ins, outs) {
    return fmt.Sprintf("%d: Ins and Outs don't Balance\n", sus.Tid), false
  } // endif not Balanced.

  if Duplicated(ins) {
    return fmt.Sprintf("%d: DoubleSpend attack: Inputs contains duplicate coins.\n", sus.Tid), false
    structures.PrintCoins(ins, "Duplicates")
  } // endif Duplicated.

  if ZeroPayment(ins) {
    return fmt.Sprintf("%d: Input contains zero valued coins.\n", sus.Tid), false
    structures.PrintCoins(ins, "Zeroes")
  } // endif Duplicated.

  uttxs, _ := InUtxos(sus, emmone)
  if len(uttxs) > len(ins) {
    return fmt.Sprintf("%d: Double Spent coins in the input.\n", sus.Tid), false
  } else if len(uttxs) < len(ins) {
    return fmt.Sprintf("%d: Counterfeit coins in the input.\n", sus.Tid), false
  } else {
    return fmt.Sprintf("%d: All good coins in the input.\n", sus.Tid), true
  } // endif not in utxos.

  enutxs := Embezzled(sus, uttxs)
  if len(enutxs) > 0 {
    structures.PrintCoins(enutxs, "Stolen?")
    return fmt.Sprintf("%d: Some Coins Not in Utxos: Counterfeit.\n", sus.Tid), false
  } // endif not in utxos.

  // It's in utxos, it balances, it's yours, it verifies.
  return fmt.Sprintf("%d: Good Tx.", sus.Tid), true

} // end func LegitimateTransaction.


// Coins out must balance coins in.
func Balanced(inns []structures.Coin, outts []structures.Coin) bool {

  insum := structures.CoinValues(inns)
  outsum := structures.CoinValues(outts)
  return insum == outsum

} // endif Balanced.


// Check the inputs list for duplicates.
func Duplicated(inn []structures.Coin) bool {

  if len(inn) == 1 { return false }

  if len(inn) == 2 {
    return ((inn[0].Cid == inn[1].Cid) && (inn[0].Seq == inn[1].Seq))

  for ix := 0; ix < len(inn)-1; ix++ {
    for jx := 1; jx < len(inn); jx++ {
      if ((inn[ix].Cid == inn[jx].Cid) && (inn[ix].Seq == inn[jx].Seq)) { return true }
    } // end for comparators. 
  } // end for all inputs.
  } // end compare two.

  return false
 
}

// Check the inputs list for zero values
func ZeroPayment(inn []structures.Coin) bool {

  for jx := 0; jx < len(inn); jx++ {
    xerox, _ := strconv.Atoi(inn[jx].Denom)
    if xerox == 0 { return true }
  } // end for zero check.

  return false
 
}

// Separate out coins and non-coins and return them.
func InUtxos(suspect structures.Transaction, emmony []structures.Coin) ([]structures.Coin, []structures.Coin) {

  otxos := []structures.Coin{}
  notxos := []structures.Coin{}
  tfrees := suspect.Inputs

  structures.PrintCoins(tfrees, "SUSPECTINPUTS")

  for ix := 0; ix < len(tfrees); ix++ {
    suscoin := tfrees[ix]
    for sx := 0; sx < len(emmony); sx++ {
      basemoney := emmony[sx]
      if (suscoin.Cid == basemoney.Cid && suscoin.Seq == basemoney.Seq) {
        otxos = append(otxos, basemoney)
      } else {
        notxos = append(notxos, basemoney)
      } // endif IdSeq match.
    } // end for M1.
  } // end for Inputs.

  return otxos, notxos

} // endif InUtxos.


// Separate good coins that are in Utxos into mine and not mine.
func Embezzled(susp structures.Transaction, utter []structures.Coin) []structures.Coin {

  yeses := 0
  notxos := []structures.Coin{}
  bezzlers := susp.Inputs

  for ix := 0; ix < len(bezzlers); ix++ {
    suscoin := bezzlers[ix]
    sushash := himitsu.DERToHash(suscoin.Owner)
    for sx := 0; sx < len(utter); sx++ {
      basem := utter[sx]

      if sushash == basem.Owner {
        yeses += 1
      } else {
        notxos = append(notxos, suscoin)
      } // endif sus in M1.

    } // end for tfrees.
  } // end for M1.

  return notxos

} // endif Embezzled.


// NotDoubleSpent: Check the ledger to see if the current transaction
// is a double spend attack.
// Get the coinbase, see if paycoins transferor is coin owner.
func NotDoubleSpent(emmone []structures.Coin, suspect structures.Transaction) bool {

  tfrees := suspect.Inputs
  result := false

  for ix := 0; ix < len(emmone); ix++ {
    basemoney := emmone[ix]
    for sx := 0; sx < len(tfrees); sx++ {
      suscoin := tfrees[sx]
      if suscoin.Cid == basemoney.Cid {
        if suscoin.Seq == basemoney.Seq {
          // Incoming Tx has full pubkey. Hash it.
          sushash := himitsu.DERToHash(suscoin.Owner)
          if sushash == basemoney.Owner {
            // fmt.Printf("NDS: Successful Payment from: %s\n", sushash)
            result = true
          } else {
            // it's a double spend.
            fmt.Printf("NDS: REJECTED: M1 owner=%s, PayCoins In says:%s\n", basemoney.Owner[:5], sushash[:5])
            result = false
          } // endif Owner.
        } // endif Seq,
      } // end Cid.

    } // end for tfrees.
  } // end for m1.

  return result

} // end func DoubleSpent.


// Bankeronly:  Only the central banker is allowed to create coins.
func BankerOnly( conf map[string]string, checking structures.Transaction, verby bool) bool {

  bankerkey := himitsu.BaseDER(conf["pubkey"])
  checkerkey := checking.Inputs[0].Owner
  return bankerkey == checkerkey

} // end func BankerOnly.

// File existence test.
func NoFile(filename string) bool {

  if _, err := os.Stat(filename); os.IsNotExist(err) {
    return true
  } else {
    return false
  } // endif stat.

} // end func NoFile.

