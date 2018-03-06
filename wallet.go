package wallet

import (
	"crypto/elliptic"
    "crypto/ecdsa"
	"crypto/sha256"
	// "fmt"
	"log"
	"strings"
	"math/big"

	"golang.org/x/crypto/ripemd160"
	
)

const version = byte(0x00)
const addressChecksumLen = 4

// Wallet stores private and public keys
type Wallet struct {
	PrivateKey []byte
	PublicKey  []byte
}

type PublicKey struct {
	elliptic.Curve
	X, Y *big.Int
}

// PrivateKey represents a ECDSA private key.
type PrivateKey struct {
	PublicKey
	D *big.Int
}

// NewWallet creates and returns a Wallet
func NewWallet() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet{private, public}

	return &wallet
}

func SetWallet(prvk string) *Wallet {
	private, public := keyPair(prvk)
	wallet := Wallet{private, public}

	return &wallet
}

// GetAddress returns wallet address
func (w Wallet) GetAddress() (address string) {
	/* See https://en.bitcoin.it/wiki/Technical_background_of_Bitcoin_addresses */

	/* Convert the public key to bytes */
	pub_bytes := w.PublicKey

	/* SHA256 Hash */
	// fmt.Println("2 - Perform SHA-256 hashing on the public key")
	sha256_h := sha256.New()
	sha256_h.Reset()
	sha256_h.Write(pub_bytes)
	pub_hash_1 := sha256_h.Sum(nil)
	// fmt.Println(ByteString(pub_hash_1))
	// fmt.Println("=======================")

	/* RIPEMD-160 Hash */
	// fmt.Println("3 - Perform RIPEMD-160 hashing on the result of SHA-256")
	ripemd160_h := ripemd160.New()
	ripemd160_h.Reset()
	ripemd160_h.Write(pub_hash_1)
	pub_hash_2 := ripemd160_h.Sum(nil)
	// fmt.Println(ByteString(pub_hash_2))
	// fmt.Println("=======================")
	/* Convert hash bytes to base58 check encoded sequence */
	address = B58checkencode(0x00, pub_hash_2)

	return address
}

// HashPubKey hashes public key
func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

const privKeyBytesLen = 32

func newKeyPair() ([]byte, []byte) {
	curve := elliptic.P256()
	// private, err := ecdsa.GenerateKey(curve, rand.Reader)
	private, err := ecdsa.GenerateKey(curve, strings.NewReader("YYw5c5AqqWLBdgRdLbVwNGZYmsvn8yPzi6RUA1LCgSwDfe3xrRKsd"))
	if err != nil {
		log.Panic(err)
	}
	d := private.D.Bytes()
	b := make([]byte, 0, privKeyBytesLen)
	priKet := PaddedAppend(privKeyBytesLen, b, d)
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return priKet, pubKey
}

func keyPair(privStr string) ([]byte, []byte) {
	curve := elliptic.P256()
	d := new(big.Int).SetBytes(HexToByte(privStr))

	priv2 := new(PrivateKey)
	priv2.PublicKey.Curve = curve
	priv2.D = d 
    priv2.PublicKey.X, priv2.PublicKey.Y = curve.ScalarBaseMult(d.Bytes())
	pubKey := append(priv2.PublicKey.X.Bytes(), priv2.PublicKey.Y.Bytes()...)
	
	b := make([]byte, 0, privKeyBytesLen)
	priKet := PaddedAppend(privKeyBytesLen, b, d.Bytes())

	return priKet, pubKey
}

// ToWIF converts a Bitcoin private key to a Wallet Import Format string.
func ToWIF(priv []byte) (wif string) {
	/* Convert bytes to base-58 check encoded string with version 0x80 */
	wif = B58checkencode(0x80, priv)

	return wif
}

