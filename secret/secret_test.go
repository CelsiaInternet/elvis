package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

// Vector producido con el codigo de migration-secret.go de catalog y send (2026-11-17),
// llave "llave-de-compatibilidad", texto "clave-de-origen".
const legacyVector = "v1:/1CUaXtbbvf7FRA1IMbxnga/NpqRDQ8GPPX+bOR0cMrEZ4vx+SSHFqSF6A=="

/**
* legacyEncrypt: Copia de encryptMigrationSecret de catalog y send antes de usar este paquete.
* @param t *testing.T, key, plain string
* @return string
**/
func legacyEncrypt(t *testing.T, key, plain string) string {
	sum := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(sum[:])
	if err != nil {
		t.Fatal(err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		t.Fatal(err)
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		t.Fatal(err)
	}
	sealed := aead.Seal(nonce, nonce, []byte(plain), nil)
	return "v1:" + base64.StdEncoding.EncodeToString(sealed)
}

/**
* mustNew: Cipher con la llave dada; falla el test si no se puede crear.
* @param t *testing.T, key string
* @return *Cipher
**/
func mustNew(t *testing.T, key string) *Cipher {
	c, err := New(key)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestRoundTrip(t *testing.T) {
	c := mustNew(t, "llave")
	for _, plain := range []string{"clave-de-origen", "", "p'ss \"ñ\" 🎸"} {
		stored, err := c.Encrypt(plain)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.HasPrefix(stored, Version) || (plain != "" && strings.Contains(stored, plain)) {
			t.Fatalf("cifrado: %s", stored)
		}
		got, err := c.Decrypt(stored)
		if err != nil || got != plain {
			t.Fatalf("descifrado: %q %v", got, err)
		}
	}
}

func TestNonceChanges(t *testing.T) {
	c := mustNew(t, "llave")
	a, _ := c.Encrypt("igual")
	b, _ := c.Encrypt("igual")
	if a == b {
		t.Fatal("dos cifrados iguales: el nonce no cambia")
	}
}

func TestTampered(t *testing.T) {
	c := mustNew(t, "llave")
	stored, _ := c.Encrypt("clave-de-origen")
	data, _ := base64.StdEncoding.DecodeString(strings.TrimPrefix(stored, Version))
	data[len(data)-1] ^= 0x01
	tampered := Version + base64.StdEncoding.EncodeToString(data)

	_, err := c.Decrypt(tampered)
	var invalid *InvalidError
	if !errors.Is(err, ErrInvalid) || !errors.As(err, &invalid) || errors.Unwrap(err) == nil {
		t.Fatalf("texto alterado: %v", err)
	}
	if !strings.HasPrefix(err.Error(), MSG_SECRET_INVALID+": ") {
		t.Fatalf("mensaje: %s", err)
	}
}

func TestInvalidFormat(t *testing.T) {
	c := mustNew(t, "llave")
	for _, stored := range []string{"texto-plano", "v2:AAAA", "v1:%%%", "v1:" + base64.StdEncoding.EncodeToString([]byte("corto"))} {
		if _, err := c.Decrypt(stored); err != ErrInvalid {
			t.Errorf("%q: %v", stored, err)
		}
	}
}

func TestOtherKey(t *testing.T) {
	stored, _ := mustNew(t, "llave").Encrypt("clave-de-origen")
	if _, err := mustNew(t, "cambiada").Decrypt(stored); !errors.Is(err, ErrInvalid) {
		t.Fatalf("otra llave: %v", err)
	}
}

func TestKeyRequired(t *testing.T) {
	if _, err := New(""); err != ErrKeyRequired {
		t.Fatalf("sin llave: %v", err)
	}
	if _, err := FromEnv("ELVIS_SECRET_TEST_EMPTY_A", "ELVIS_SECRET_TEST_EMPTY_B"); err != ErrKeyRequired {
		t.Fatalf("sin variables: %v", err)
	}
	if _, err := Encrypt("x", "ELVIS_SECRET_TEST_EMPTY_A"); err != ErrKeyRequired {
		t.Fatalf("Encrypt sin llave: %v", err)
	}
	if _, err := Decrypt(Version+"AAAA", "ELVIS_SECRET_TEST_EMPTY_A"); err != ErrKeyRequired {
		t.Fatalf("Decrypt sin llave: %v", err)
	}
}

func TestFromEnvFallback(t *testing.T) {
	// Nombres propios del test: envar guarda en memoria el primer valor que lee de cada variable
	os.Setenv("ELVIS_SECRET_TEST_FALLBACK", "respaldo")
	defer os.Unsetenv("ELVIS_SECRET_TEST_FALLBACK")

	if key := EnvKey("ELVIS_SECRET_TEST_PRIMARY_EMPTY", "ELVIS_SECRET_TEST_FALLBACK"); key != "respaldo" {
		t.Fatalf("respaldo: %q", key)
	}
	stored, err := Encrypt("clave", "ELVIS_SECRET_TEST_PRIMARY_EMPTY", "ELVIS_SECRET_TEST_FALLBACK")
	if err != nil {
		t.Fatal(err)
	}
	plain, err := mustNew(t, "respaldo").Decrypt(stored)
	if err != nil || plain != "clave" {
		t.Fatalf("%q %v", plain, err)
	}
	plain, err = Decrypt(stored, "ELVIS_SECRET_TEST_FALLBACK")
	if err != nil || plain != "clave" {
		t.Fatalf("Decrypt: %q %v", plain, err)
	}
}

func TestLegacyCompatibility(t *testing.T) {
	c := mustNew(t, "llave-de-compatibilidad")
	plain, err := c.Decrypt(legacyVector)
	if err != nil || plain != "clave-de-origen" {
		t.Fatalf("vector fijo: %q %v", plain, err)
	}

	stored := legacyEncrypt(t, "llave-de-compatibilidad", "otra clave")
	plain, err = c.Decrypt(stored)
	if err != nil || plain != "otra clave" {
		t.Fatalf("codigo anterior: %q %v", plain, err)
	}
}
