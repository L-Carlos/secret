package secret

import (
	"fmt"

	"github.com/L-Carlos/secret/encrypt"
)

type Vault struct {
	encodingKey string
	keyValues   map[string]string
}

func MemoryVault(encodingKey string) Vault {
	return Vault{
		encodingKey: encodingKey,
		keyValues:   map[string]string{}}
}

func (v *Vault) Get(key string) (string, error) {
	hexValue, ok := v.keyValues[key]
	if !ok {
		return "", fmt.Errorf("secret: no value for (%s) key", key)
	}
	value, err := encrypt.Decrypt(v.encodingKey, hexValue)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (v *Vault) Set(key, value string) error {
	hexValue, err := encrypt.Encrypt(v.encodingKey, value)
	if err != nil {
		return err
	}
	v.keyValues[key] = hexValue
	return nil
}
