// Package secret cifra y descifra secretos que un servicio guarda (por ejemplo la contrasena de una
// base origen) con AES-256-GCM. El texto guardado es Version + base64(nonce || cifrado) y la llave
// es sha256 del secreto configurado. Guardar o no el secreto lo decide cada servicio; el paquete
// solo cifra y descifra.
//
// A diferencia de utility.EncryptMessage (CFB con SECRET tal cual como llave), detecta un texto
// alterado y acepta una llave de cualquier largo.
package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"strings"

	"github.com/celsiainternet/elvis/envar"
)

// Version: prefijo del texto cifrado, para poder cambiar el esquema despues
const Version = "v1:"

var (
	// ErrKeyRequired: la llave esta vacia (ninguna variable de entorno tiene valor)
	ErrKeyRequired = errors.New(MSG_SECRET_KEY_REQUIRED)
	// ErrInvalid: el texto no tiene el prefijo Version, no es base64, se altero o la llave cambio
	ErrInvalid = errors.New(MSG_SECRET_INVALID)
)

/**
* InvalidError: Error de Decrypt cuando GCM no abre el texto; envuelve ErrInvalid.
**/
type InvalidError struct {
	Cause error
}

/**
* Error: Mensaje de ErrInvalid seguido del error de GCM.
* @return string
**/
func (e *InvalidError) Error() string {
	return ErrInvalid.Error() + ": " + e.Cause.Error()
}

/**
* Is: Hace que errors.Is(err, ErrInvalid) sea verdadero.
* @param target error
* @return bool
**/
func (e *InvalidError) Is(target error) bool {
	return target == ErrInvalid
}

/**
* Unwrap: Error de GCM que causo el fallo.
* @return error
**/
func (e *InvalidError) Unwrap() error {
	return e.Cause
}

/**
* Cipher: Cifrador AES-256-GCM con una llave fija. Es seguro para uso concurrente.
**/
type Cipher struct {
	aead cipher.AEAD
}

/**
* EnvKey: Valor de la primera variable de entorno no vacia (envar.GetStr). Sin nombres usa SECRET.
* @param names ...string
* @return string
**/
func EnvKey(names ...string) string {
	if len(names) == 0 {
		names = []string{"SECRET"}
	}

	for _, name := range names {
		if value := envar.GetStr("", name); value != "" {
			return value
		}
	}

	return ""
}

/**
* New: Crea un Cipher con la llave sha256(key). Sirve para inyectar la llave en tests.
* @param key string
* @return *Cipher, error (ErrKeyRequired si key esta vacia)
**/
func New(key string) (*Cipher, error) {
	if key == "" {
		return nil, ErrKeyRequired
	}

	sum := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(sum[:])
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &Cipher{aead: aead}, nil
}

/**
* FromEnv: Crea un Cipher con la primera variable no vacía (SECRET por defecto).
* @param names ...string
* @return *Cipher, error (ErrKeyRequired si ninguna tiene valor)
**/
func FromEnv(names ...string) (*Cipher, error) {
	return New(EnvKey(names...))
}

/**
* Encrypt: Cifra plain: Version + base64(nonce || cifrado), con un nonce aleatorio en cada llamada.
* @param plain string
* @return string, error
**/
func (c *Cipher) Encrypt(plain string) (string, error) {
	nonce := make([]byte, c.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	sealed := c.aead.Seal(nonce, nonce, []byte(plain), nil)
	return Version + base64.StdEncoding.EncodeToString(sealed), nil
}

/**
* Decrypt: Descifra lo que guardo Encrypt.
* @param stored string
* @return string, error (ErrInvalid sin prefijo o base64 roto; *InvalidError si se altero o la llave cambio)
**/
func (c *Cipher) Decrypt(stored string) (string, error) {
	if !strings.HasPrefix(stored, Version) {
		return "", ErrInvalid
	}

	data, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(stored, Version))
	size := c.aead.NonceSize()
	if err != nil || len(data) < size {
		return "", ErrInvalid
	}

	plain, err := c.aead.Open(nil, data[:size], data[size:], nil)
	if err != nil {
		return "", &InvalidError{Cause: err}
	}

	return string(plain), nil
}

/**
* Encrypt: Cifra plain con la llave de la primera variable no vacia de envNames (sin nombres, SECRET).
* @param plain string, envNames ...string
* @return string, error
**/
func Encrypt(plain string, envNames ...string) (string, error) {
	c, err := FromEnv(envNames...)
	if err != nil {
		return "", err
	}

	return c.Encrypt(plain)
}

/**
* Decrypt: Descifra stored con la llave de la primera variable no vacia de envNames (sin nombres, SECRET).
* @param stored string, envNames ...string
* @return string, error
**/
func Decrypt(stored string, envNames ...string) (string, error) {
	c, err := FromEnv(envNames...)
	if err != nil {
		return "", err
	}

	return c.Decrypt(stored)
}
