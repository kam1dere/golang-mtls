package cipher_test

import (
	"mtls/pkg/cipher"
	"testing"

	"mtls/pkg/curve25519"
	"mtls/pkg/ed25519"
)

func provideCipher() (c *cipher.Cipher, err error) {
	pub, priv, err := ed25519.GenerateKeyPair()
	if err != nil {
		return nil, err
	}

	var shared []byte

	shared, err = curve25519.GenerateSharedKey(pub, priv)
	if err != nil {
		return nil, err
	}

	c, err = cipher.NewCipher(shared)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func TestCipher(t *testing.T) {
	t.Parallel()

	testCipher, err := provideCipher()
	if err != nil {
		t.Fatal(err.Error())
	}

	origin := []byte("a special secret message")

	var encoded []byte

	encoded, err = testCipher.Encode(origin)
	if err != nil {
		t.Fatal(err.Error())
	}

	var decoded []byte

	decoded, err = testCipher.Decode(encoded)
	if err != nil {
		t.Fatal(err.Error())
	}

	if string(origin) != string(decoded) {
		t.Fatal("Origin and Decoded are not equal")
	}
}

func BenchmarkCipherEncode(b *testing.B) {
	testCipher, err := provideCipher()
	if err != nil {
		b.Fatal(err.Error())
	}

	origin := []byte("a special secret message")

	for range b.N {
		_, err = testCipher.Encode(origin)
		if err != nil {
			b.Fatal(err.Error())
		}
	}
}

func BenchmarkCipherDecode(b *testing.B) {
	origin := []byte("a special secret message")

	testCipher, err := provideCipher()
	if err != nil {
		b.Fatal(err.Error())
	}

	var encoded []byte

	encoded, err = testCipher.Encode(origin)
	if err != nil {
		b.Fatal(err.Error())
	}

	for range b.N {
		_, err = testCipher.Decode(encoded)
		if err != nil {
			b.Fatal(err.Error())
		}
	}
}
