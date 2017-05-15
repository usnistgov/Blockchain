package main

import (
  "fmt"
)

func main() {
  sup := 64; bens := 6

  q, r := DivRem(sup, bens)
  fmt.Printf("%d/%d = %d, %d %% %d = %d\n", sup, bens, q, sup, bens, r)

} // end func main.

func DivRem(supply int, bennies int) (int, int) {

  q := supply / bennies
  r := supply % bennies

  return q, r

} // end func DivRem.

