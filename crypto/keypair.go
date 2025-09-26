package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"

	"github.com/772005himanshu/Mingo-Blockchain/types"
)


type PubKey []byte  // Public is the slice of the bytes



type PrivateKey struct {
	key *ecdsa.PrivateKey
}

func (k PrivateKey) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, k.key, data)
	if err != nil {
		return nil, err
	}

	return &Signature{
		R: r,
		S: s,
	}, nil // firts Mistake Here
}

func GeneratePrivateKey() PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader) // Never Use the Rand due to determinitic Nature of the Blockchain
	if err != nil {
		panic(err)
	}

	return PrivateKey{
		key: key,
	}
}

func (k PrivateKey) PublicKey() PubKey {
	return elliptic.MarshalCompressed(k.key.PublicKey, k.key.PublicKey.X, k.key.PublicKey.Y)  // From This we can get the Compressed Verion of the Public key 
}

type PublicKey []byte


func (k PublicKey) Address() types.Address {
	h := sha256.Sum256(k)  // We can Directly use the k because this is in the bytes

	return types.AddressFromBytes(h[len(h)-20:])
}

type Signature struct {
	S, R *big.Int
}

func (sig Signature) String() string {
	// buf := new(bytes.Buffer)
	// buf.Write(sig.S.Bytes())
	// buf.Write(sig.R.Bytes())

	// return buf.String()  > Check we can do this or not ?

	b := append(sig.S.Bytes(), sig.R.Bytes()...)  // ... append usually takes values one by one // If we want to append the Complete slice so we use ...
	return hex.EncodeToString(b)
}

func (sig Signature) Verify(pubKey PubKey, data []byte) bool {
	x, y := elliptic.UnmarshalCompressed(elliptic.P256(),pubKey)
	key := &ecdsa.PublicKey{
		Curve : elliptic.P256(),
		X: x,
		Y: y,
	}  // From here we get the slice of publicKey that is in the PublicKey

	// now we want to convert to the Ecdsa Pubkey , so we pubt it in the ecdsa verify 

	return ecdsa.Verify(key, data, sig.R, sig.S)
}
