package secret

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/L-Carlos/secret/encrypt"
)

type Vault struct {
	encodingKey string
	filepath    string
	mutex       sync.Mutex
	keyValues   map[string]string
}

func FileVault(encodingKey, filePath string) *Vault {
	return &Vault{
		encodingKey: encodingKey,
		filepath:    filePath,
		keyValues:   map[string]string{},
	}

}

func (v *Vault) loadKeyValues() error {
	f, err := os.Open(v.filepath)
	if err != nil {
		if os.IsNotExist(err) {
			v.keyValues = map[string]string{}
			return nil
		} else {
			return err
		}
	}
	defer f.Close()

	var sb strings.Builder
	_, err = io.Copy(&sb, f)
	if err != nil {
		return err
	}

	decryptedJSON, err := encrypt.Decrypt(v.encodingKey, sb.String())
	if err != nil {
		return err
	}

	reader := strings.NewReader(decryptedJSON)
	dec := json.NewDecoder(reader)

	return dec.Decode(&v.keyValues)
}

func (v *Vault) saveKeyValues() error {
	var sb strings.Builder
	enc := json.NewEncoder(&sb)
	err := enc.Encode(v.keyValues)
	if err != nil {
		return err
	}

	encryptedJSON, err := encrypt.Encrypt(v.encodingKey, sb.String())
	if err != nil {
		return err
	}

	f, err := os.OpenFile(v.filepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = io.WriteString(f, encryptedJSON)

	return err
}

func (v *Vault) Get(key string) (string, error) {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	err := v.loadKeyValues()
	if err != nil {
		return "", err
	}
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
	v.mutex.Lock()
	defer v.mutex.Unlock()

	hexValue, err := encrypt.Encrypt(v.encodingKey, value)
	if err != nil {
		return err
	}

	err = v.loadKeyValues()
	if err != nil {
		return err
	}

	v.keyValues[key] = hexValue

	return v.saveKeyValues()
}
