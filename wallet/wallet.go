package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"

	"golang.org/x/crypto/ripemd160"
)

// NOTE :: ripemd160 is deprecated and is not recommended to use in moder cryptography due to potential vulnerabilies, therefore we will rather try to truncate sha-256 to 160 bits

const (
	checksumLength = 4
	version = byte(0x00)
)

type Wallet struct{
	PrivateKey ecdsa.PrivateKey
	PublicKey []byte
}

func (w Wallet) Address() []byte{
	pubHash := PublicKeyHash(w.PublicKey)

	versionedHash := append([]byte{version},pubHash...)
	checksum := Checksum(versionedHash)

	fullHash := append(versionedHash,checksum...)
	address := Base58Encode(fullHash)

	return address
}

func ValidateAddress(address string) bool{
	pubKeyHash := Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-checksumLength:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-checksumLength]
	targetChecksum := Checksum(append([]byte{version},pubKeyHash...))

	return bytes.Equal(actualChecksum,targetChecksum)
}

func NewKeyPair() (ecdsa.PrivateKey, []byte){
	curve := elliptic.P256()

	private,err := ecdsa.GenerateKey(curve,rand.Reader)
	Handle(err)

	pub := append(private.PublicKey.X.Bytes(),private.PublicKey.Y.Bytes()...)
	return *private,pub
}

func (w Wallet) MarshalJSON() ([]byte, error) {
	mapStringAny := map[string]any{
		"PrivateKey": map[string]any{
			"D": w.PrivateKey.D,
			"PublicKey": map[string]any{
				"X": w.PrivateKey.PublicKey.X,
				"Y": w.PrivateKey.PublicKey.Y,
			},
			"X": w.PrivateKey.X,
			"Y": w.PrivateKey.Y,
		},
		"PublicKey": w.PublicKey,
	}
	return json.Marshal(mapStringAny)
}

func MakeWallet() *Wallet{
	private, public := NewKeyPair()
	wallet := Wallet{private, public}

	return &wallet
}

func PublicKeyHash(pubKey []byte) []byte{
	pubHash := sha256.Sum256(pubKey)

	hasher := ripemd160.New()
	_,err := hasher.Write(pubHash[:])
	Handle(err)

	publicRipMD := hasher.Sum(nil)

	// truncatedHash := pubHash[:20]

	return publicRipMD
}

func Checksum(payload []byte) []byte{
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checksumLength]
}