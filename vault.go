package secret

import "fmt"

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
	if value, ok := v.keyValues[key]; ok {
		return value, nil
	}
	return "", fmt.Errorf("secret: no value for (%s) key", key)
}

func (v *Vault) Set(key, value string) {
	v.keyValues[key] = value
}
