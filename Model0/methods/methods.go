package methods

import (
  "bufio"
  "encoding/base64"
  "encoding/json"
  "fmt"
  "log"
  "currency/himitsu"
  "currency/structures"
  "io/ioutil"
  "os"
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
  return mapvec

} // getConfigs.

// List the files in a given directory, with a given suffix.
func GetDir(dirname string, suffix string) (fileslice []string) {
  var fullpath string
  files, err := ioutil.ReadDir(dirname)
  if err != nil {
    log.Fatal(err)
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
}

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


// Count Banker coins as 'M0'.
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
func M1(led []structures.Transaction) []structures.Coin {
  coincount := []structures.Coin{}
  originals, assignees := 0, 0

  for ix := 0; ix < len(led); ix++ {
    nextx := led[ix]

    switch(nextx.Ttyp) {

      // No coins involved.
      case "BigBang":
        continue

      // Add to original coin count.
      case "CreateCoins":
        coincount = CountCoins(nextx.Outputs, coincount)
        originals = len(coincount)
        fmt.Printf("M1: %d coins up to transaction %d.\n", originals, nextx.Tid) 

      // No change to coin count. Change coin owners.
      case "PayCoins":
        coincount = ReassignedCoins(nextx.Outputs, coincount)
        assignees = len(coincount)

      // No other Trx types expected.
      default:
        fmt.Printf("%s: Bad Transaction Type.\n", nextx.Ttyp)
        os.Exit(1)

    } // end switch Ttyp.

  } // end for led.

  if originals != assignees {
    fmt.Printf("Warning: M1: %d coins created, %d coins assigned.\n", originals, assignees) 
  } // endif noriginals.

  return coincount

} // end func M1.


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
  CheckErrorInst(1, errb, false)
  verifier, errv := himitsu.ParseDERKey(pubb)
  CheckErrorInst(2, errv, true)
  signet, errs := json.Marshal(onetran.Outputs)
  CheckErrorInst(3, errs, false)
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
func ProcessNewTransactions(config map[string]string, lodger []structures.Transaction) ([]structures.Transaction, []structures.Transaction) {

  newin := []structures.Transaction{}
  newout := []structures.Transaction{}
  rejects := []structures.Transaction{}

  if NoFile(config["newtxs"]) {
    fmt.Println("No New Transactions.")
    return newout, rejects
  } // endif NoFile.

  newin = LoadLedger(config["newtxs"], newin)
  emmy := M1(lodger)  // existing coinbase to check for double spends.
  fmt.Printf("%d New Transactions in.\n", len(newin))

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
        if NotDoubleSpent(emmy, checking, false) {
          newout = append(newout, checking)
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
  return newout, rejects

} // end func ProcessNewTransactions.


func OLDProcessNewTransactions(config map[string]string, lodger []structures.Transaction) ([]structures.Transaction, []structures.Transaction) {

  rejects := []structures.Transaction{}
  newones := []structures.Transaction{}
  verby := false

  if config["verbose"] == "true" {
    verby = true
  } // endif verby.

  if NoFile(config["newtxs"]) {
    fmt.Println("No New Transactions.")
    return lodger, rejects
  } // endif NoFile.

  newones = LoadLedger(config["newtxs"], newones)
  emmy := M1(lodger)
  result := true

  // verify  and check for double spends.
  fmt.Printf("Checking %d new transactions for double spend.\n", len(newones))
  for ix := 0; ix < len(newones); ix++ {
    checking := newones[ix]
    structures.PrintTransaction(checking, "\nNEW TRANSACTION:\n")
    if VerifyTransaction(checking, verby) != nil {
      result = false
    } else {
      fmt.Printf("Transaction %d verified: %s\n", checking.Tid, checking.Ttyp)
      switch(checking.Ttyp) {
        case "CreateCoins":
          result = BankerOnly(config, checking, verby)
          if !result {
            fmt.Println("REJECTED: only Banker can Create Coinage.")
          } // endif reult.

        case "PayCoins":
          result = NotDoubleSpent(emmy, checking, verby)

        default:
          result = false

      } // end switch Ttyp.

      if result {
        fmt.Printf("PNTs: Len lodger before %d, len m1 before %d\n", len(lodger), len(emmy))
        lodger = append(lodger, checking)
        emmy = M1(lodger)
        structures.PrintCoins(emmy, "OwnerOut")
        fmt.Printf("PNTs: Len lodger after %d, len m1 after %d\n", len(lodger), len(emmy))
      } else {
        rejects = append(rejects, checking)
        fmt.Printf("PNTs: Len rejects after %d\n", len(rejects))
      } // endif result.

    } // endif VerifyTransaction.

  } // end for newones.

  // Delete new transactions file.
  err := os.Remove(config["newtxs"])
  CheckError(err, true)

  fmt.Printf("ProcessNewTxs: Returning %d ledger, %d rejects\n", len(lodger), len(rejects))
  return lodger, rejects

} // end func ProcessNewTransactions.


// DoubleSpent: Check the ledger to see if the current transaction
// is a double spend attack.
// Get the coinbase, see if paycoins transferor is coin owner.
func NotDoubleSpent(emmone []structures.Coin, suspect structures.Transaction, verbose bool) bool {

  if verbose {
    structures.PrintTransaction(suspect, suspect.Ttyp)
  } // endif verbose.

  tfrees := suspect.Inputs
  result := false

  fmt.Printf("\n\nDoubleSpent: coin count %d, Input coins are %d\n", len(emmone), len(tfrees))
  for ix := 0; ix < len(emmone); ix++ {
    basemoney := emmone[ix]
    // fmt.Println("Basemoney Coin:")
    // structures.PrintCoin(basemoney)
    for sx := 0; sx < len(tfrees); sx++ {
      suscoin := tfrees[sx]
      // fmt.Println("New Transaction Coin:")
      // structures.PrintCoin(suscoin)
      if suscoin.Cid == basemoney.Cid {
        if suscoin.Seq == basemoney.Seq {
          // Incoming Tx has full pubkey. Hash it.
          sushash := himitsu.DERToHash(suscoin.Owner)
          if sushash == basemoney.Owner {
            fmt.Printf("DoubleSpent: Successful Payment from: %s\n", sushash)
            result = true
          } else {
            // it's a double spend.
            fmt.Println("REJECTED: Different Owner: Double Spend Attempt.")
            fmt.Println("M1: ", basemoney.Owner)
            fmt.Println("In: ", sushash)
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

