package himitsu

import (
  "crypto"
  "crypto/rand"
  "crypto/rsa"
  "crypto/sha256"
  "crypto/x509"
  "encoding/base64"
  "encoding/json"
  "encoding/pem"
  "errors"
  "fmt"
  "io/ioutil"
  "os"
  "strings"
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


// Convert the input byte string to Base64.
func Base64Transaction(tx []byte) (bs string) {
  jtx, _ := json.Marshal(tx)
  bs = base64.StdEncoding.EncodeToString(jtx)
  return bs
} // end Base64Transaction.


// Get Base64 of DER form of public key.
func BaseDER(path string) string {
  dat, _ := ioutil.ReadFile(path)
  block, _ := pem.Decode(dat)
  if block == nil {
    return "No Key Found."
  } // end if.

  pubout := base64.StdEncoding.EncodeToString(block.Bytes)
  return pubout

} // end func BaseDER.

// Convert the DER form publick key into its Hash.
func DERToHash(der string) string {

  bites, _ := base64.StdEncoding.DecodeString(der)
  
  return Bashit(bites)

} // end func DERToHash.

// loadPublicKey loads and parses a PEM encoded private key file.
func LoadPublicKey(path string) (Unsigner, error) {
  dat, _ := ioutil.ReadFile(path)
  return ParsePublicKey([]byte(dat))
}

// DisplayPublicKey returns the PEM encoded public key file.
func DisplayPublicKey(path string) (string) {
  dat, _ := ioutil.ReadFile(path)
  return string(dat)
}

// HashPublicKey returns the base64 encoded hash of the DER form public key.
func HashPublicKey(path string) (string) {
  // fmt.Printf("Opening %s in HashPublicKey ...\n", path)
  dat, _ := ioutil.ReadFile(path)
  block, _ := pem.Decode([]byte(dat))
  hashed :=  Bashit(block.Bytes)
  return hashed
}

// parsePublicKey parses a PEM encoded public key.
func ParsePublicKey(pemBytes []byte) (Unsigner, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("ssh: no key found")
	}

	var rawkey interface{}
	switch block.Type {
	case "PUBLIC KEY":
		rsa, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rawkey = rsa
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %q", block.Type)
	}

	return NewUnsignerFromKey(rawkey)
}

// ParseDERKey parses a DER form public key.
func ParseDERKey(block []byte) (Unsigner, error) {

	var rawkey interface{}
	rsa, err := x509.ParsePKIXPublicKey(block)
	if err != nil {
		return nil, err
	}
	rawkey = rsa

	return NewUnsignerFromKey(rawkey)
} // end ParseDERKey.


// loadPrivateKey loads and parses a PEM encoded private key file.
func LoadPrivateKey(path string) (Signer, error) {
  dat, _ := ioutil.ReadFile(path)
  if !strings.Contains(string(dat), "PRIVATE") {
    fmt.Println(path, " MUST contain the Private Key.")
    os.Exit(0)
  } // if not PRIVATE.

  return ParsePrivateKey([]byte(dat))
}

// parsePrivateKey parses a PEM encoded private key.
func ParsePrivateKey(pemBytes []byte) (Signer, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("ssh: no key found")
	}

	var rawkey interface{}
	switch block.Type {
	case "RSA PRIVATE KEY":
		rsa, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rawkey = rsa
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %q", block.Type)
	}
	return NewSignerFromKey(rawkey)
}

// A Signer can create signatures that verify against a public key.
type Signer interface {
	// Sign returns raw signature for the given data. This method
	// will apply the hash specified for the keytype to the data.
	Sign(data []byte) ([]byte, error)
}

// A Signer can create signatures that verify against a public key.
type Unsigner interface {
	// Sign returns raw signature for the given data. This method
	// will apply the hash specified for the keytype to the data.
	Unsign(data[]byte, sig []byte) error
}

func NewSignerFromKey(k interface{}) (Signer, error) {
	var sshKey Signer
	switch t := k.(type) {
	case *rsa.PrivateKey:
		sshKey = &rsaPrivateKey{t}
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %T", k)
	}
	return sshKey, nil
}

func NewUnsignerFromKey(k interface{}) (Unsigner, error) {
	var sshKey Unsigner
	switch t := k.(type) {
	case *rsa.PublicKey:
		sshKey = &rsaPublicKey{t}
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %T", k)
	}
	return sshKey, nil
}

type rsaPublicKey struct {
	*rsa.PublicKey
}

type rsaPrivateKey struct {
	*rsa.PrivateKey
}

// Sign signs data with rsa-sha256
func (r *rsaPrivateKey) Sign(data []byte) ([]byte, error) {
	h := sha256.New()
	h.Write(data)
	d := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, r.PrivateKey, crypto.SHA256, d)
}

// Unsign verifies the message using a rsa-sha256 signature
func (r *rsaPublicKey) Unsign(message []byte, sig []byte) error {
	h := sha256.New()
	h.Write(message)
	d := h.Sum(nil)
	return rsa.VerifyPKCS1v15(r.PublicKey, crypto.SHA256, d, sig)
}


