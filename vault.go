package secret

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
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

func (v *Vault) writeKeyValues(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(v.keyValues)
}

func (v *Vault) readKeyValues(r io.Reader) error {
	dec := json.NewDecoder(r)
	return dec.Decode(&v.keyValues)
}

func (v *Vault) save() error {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	f, err := os.OpenFile(v.filepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer f.Close()

	w, err := encrypt.EncryptWriter(v.encodingKey, f)
	if err != nil {
		return err
	}

	return v.writeKeyValues(w)
}

func (v *Vault) load() error {
	v.mutex.Lock()
	defer v.mutex.Unlock()

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

	r, err := encrypt.DecryptReader(v.encodingKey, f)
	if err != nil {
		return err
	}

	return v.readKeyValues(r)
}

func (v *Vault) Get(key string) (string, error) {
	err := v.load()
	if err != nil {
		return "", err
	}
	value, ok := v.keyValues[key]
	if !ok {
		return "", fmt.Errorf("secret: no value for (%s) key", key)
	}
	return value, nil
}

func (v *Vault) Set(key, value string) error {
	err := v.load()
	if err != nil {
		return err
	}

	v.keyValues[key] = value

	return v.save()
}

func (v *Vault) List() error {
	err := v.load()
	if err != nil {
		return err
	}

	for k, v := range v.keyValues {
		fmt.Printf("%s=%s\n", k, v)
	}

	return nil
}
