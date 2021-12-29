package cobra

import (
	"fmt"

	"github.com/L-Carlos/secret"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "gets a secret from your secret storage",
	Run: func(cmd *cobra.Command, args []string) {
		v := secret.FileVault(encodingKey, secretsFile())
		key := args[0]
		value, err := v.Get(key)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%s=%s\n", key, value)
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
}
