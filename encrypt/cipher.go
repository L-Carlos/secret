package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
)

func newCipherBlock(key string) (cipher.Block, error) {
	hasher := sha256.New()
	io.WriteString(hasher, key)
	cipherKey := hasher.Sum(nil)
	return aes.NewCipher(cipherKey)
}

func encryptStream(key string, iv []byte) (cipher.Stream, error) {
	block, err := newCipherBlock(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewCFBEncrypter(block, iv), nil
}

func decryptStream(key string, iv []byte) (cipher.Stream, error) {
	block, err := newCipherBlock(key)
	if err != nil {
		return nil, err
	}

	return cipher.NewCFBDecrypter(block, iv), nil
}

// EncryptWriter returns a writer that will write encrypted data
// to the original writer
func EncryptWriter(key string, w io.Writer) (*cipher.StreamWriter, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream, err := encryptStream(key, iv)
	if err != nil {
		return nil, err
	}

	n, err := w.Write(iv)
	if n < len(iv) || err != nil {
		return nil, fmt.Errorf("encrypt: unable to write full iv to writer")
	}

	return &cipher.StreamWriter{S: stream, W: w}, nil
}

// Encrypt takes a key and plaintext and return an hex representation
// of the encrypted value.
func Encrypt(key, plaintext string) (string, error) {
	cipherText := make([]byte, aes.BlockSize+len(plaintext))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream, err := encryptStream(key, iv)
	if err != nil {
		return "", err
	}

	stream.XORKeyStream(cipherText[aes.BlockSize:], []byte(plaintext))

	return fmt.Sprintf("%x", cipherText), nil
}

// Decrypt takes a key and a chipherHex (hexed representation
// of the ciphertext) and decrypts it.
func Decrypt(key, cipherHex string) (string, error) {
	cipherText, err := hex.DecodeString(cipherHex)
	if err != nil {
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		return "", fmt.Errorf("encrypt: cipher too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream, err := decryptStream(key, iv)
	if err != nil {
		return "", err
	}

	// XORKeyStream can work in-place if the two arguments are the same
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}

// DecryptReader returns a reader that will decrypt data
// from the provided reader
func DecryptReader(key string, r io.Reader) (*cipher.StreamReader, error) {
	iv := make([]byte, aes.BlockSize)
	n, err := r.Read(iv)
	if n < len(iv) || err != nil {
		return nil, fmt.Errorf("encrypt: unable to read the full iv from reader")
	}

	stream, err := decryptStream(key, iv)
	if err != nil {
		return nil, err
	}

	return &cipher.StreamReader{S: stream, R: r}, nil
}
