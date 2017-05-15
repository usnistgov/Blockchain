package main


import (
  "flag"
  "fmt"
  "os"
  "time"
)

var ballons = map[string]int{}
var baltime = map[string]int64{}
var rollover = 60000

// Print all Transactions in the ledger or newtxes.
func main() {

  flag.Parse()
  flay := flag.Args()
  if len(flay) != 1 {
    fmt.Println("Usage: twomaps emptyarg")
    os.Exit(1)
  } // endif flay.

  key := "wVwP"
  ballons[key] = 37
  baltime[key] = int64(time.Now().UnixNano())
  fmt.Printf("Balance for %s = %d at Time: %d\n", key, ballons[key], baltime[key])
  time.Sleep(time.Second)
  if int64(time.Now().UnixNano()) - baltime[key] > int64(rollover) {
    baltime[key] = int64(time.Now().UnixNano())
  } // endif rollover.

  fmt.Printf("Balance for %s = %d at Time: %d\n", key, ballons[key], baltime[key])

} // end main.


