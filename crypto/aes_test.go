package crypto

import (
	"fmt"
	"testing"
)

func TestAesExample(t *testing.T) {
	key := []byte("secretkey1234567")
	plaintext := []byte("Hello, world!")

	ciphertext, err := Encrypt(key, plaintext)
	if err != nil {
		panic(err)
	}

	decrypted, err := Decrypt(key, ciphertext)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", decrypted)
}
