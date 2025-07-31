package crypto

import (
	"testing"
	"fmt"
)

func TestGeneratePrivateKey(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.PublicKey()
	address := pubKey.Address()

	fmt.Println(address)
}