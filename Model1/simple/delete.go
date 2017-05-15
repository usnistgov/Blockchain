package main

import (
  "fmt"
)

func main() {
  var unspent = []int{3}
  fmt.Println(unspent)
  unspent = append(unspent[:0], unspent[1:]...)
  fmt.Println(unspent)

} // end func main.

