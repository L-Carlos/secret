package cobra

import (
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "secret",
	Short: "Secret is an API key and other secrets manager",
}

var encodingKey string

func init() {
	RootCmd.PersistentFlags().StringVarP(&encodingKey, "key", "K", "",
		"the key to use when encoding and decoding secrets")
}

func secretsFile() string {
	home, _ := homedir.Dir()
	return filepath.Join(home, ".secret")
}
