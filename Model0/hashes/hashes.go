package hashes

import (
  "crypto/sha256"
  "encoding/base64"
  "encoding/json"
)


// Hash the input string and return base64 encoded.
func Hashit(tox string) string {
  h:= sha256.New()
  h.Write([]byte(tox))
  bs := h.Sum([]byte{})
  str := base64.StdEncoding.EncodeToString(bs)
  return str
} // end Hashit.

// Hash the input string and return as bin.
func Hashbin(tox string) []byte {
  h:= sha256.New()
  h.Write([]byte(tox))
  bs := h.Sum([]byte{})
  return bs 
} // end Hashbin.


// Hash the input byte string and return base64 encoded.
func Bashit(tox []byte) string {
  h:= sha256.New()
  h.Write([]byte(tox))
  bs := h.Sum([]byte{})
  str := base64.StdEncoding.EncodeToString(bs)
  return str
} // end Bashit.

// Hash the input byte string and return as bin.
func Bashbin(tox []byte) []byte {
  h:= sha256.New()
  h.Write([]byte(tox))
  bs := h.Sum([]byte{})
  return bs 
} // end Bashbin.


// Convert the input Transaction to Base64.
func Base64Transaction(tx []byte) (bs string) {
  jtx, _ := json.Marshal(tx)
  bs = base64.StdEncoding.EncodeToString(jtx)
  return bs
} // end Base64Transaction.

