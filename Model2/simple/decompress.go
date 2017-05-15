package main

import (
  "compress/flate"
  "flag"
  "fmt"
  "io"
  "os"
)

func main() {

  flag.Parse()
  flay := flag.Args()
  if len(flay) != 2 {
    fmt.Println("Usage: decompress <smallin> <bigout>")
    os.Exit(1)
  } // endif flay.

        decompress(flay[0], flay[1])
}

func decompress(inputFile, outputFile string) {
        i, _ := os.Open(inputFile)
        defer i.Close()
        f := flate.NewReader(i)
        defer f.Close()
        o, _ := os.Create(outputFile)
        defer o.Close()
        io.Copy(o, f)
}

