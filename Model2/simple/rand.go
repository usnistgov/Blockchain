package main

import (
  "fmt"
  "math/rand"
  "strings"
  "time"
)

var row, col = 0, 0
var pairs = [][]string{
{"alice/bob", "bob/chaz", "chaz/dave", "dave/ellie", "ellie/fiona", "fiona/alice"}, 
{"alice/chaz", "bob/dave", "chaz/ellie", "dave/fiona", "ellie/alice", "fiona/bob"}, 
{"alice/dave", "bob/ellie", "chaz/fiona", "dave/alice", "ellie/bob", "fiona/chaz"}, 
{"alice/ellie", "bob/fiona", "chaz/alice", "dave/bob", "ellie/chaz", "fiona/dave"}, 
{"alice/fiona", "bob/alice", "chaz/bob", "dave/chaz", "ellie/dave", "fiona/ellie"}, 
 } // endpairs.

func Autoconstructor(span []int) []string {
  cmds := []string{}

  rand.Seed(int64(time.Now().Unix()))
  for ix := 0; ix < 120; ix++ {
    randone := Randy(span)
    sender, receiver := ReturnPair()
    onecmd := fmt.Sprintf("PayCoins users/%s.conf %d users/%s.pub", sender, randone, receiver)
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
  sen := senrec[0]; rec := senrec[1] 
  col += 1
  if col == 6 { 
    col = 0; row += 1
    if row == 5 { row = 0 }
  }
  return sen, rec
} // end func ReturnPair.

func main() {
  randramp := []int{1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20}
  commands := Autoconstructor(randramp)
  for ix := 0; ix < len(commands); ix++ {
    fmt.Println(commands[ix])
  } // end for commands.
} // end func main.

