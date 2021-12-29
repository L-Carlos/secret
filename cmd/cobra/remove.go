package cobra

import (
	"fmt"

	"github.com/L-Carlos/secret"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "removes a secret from your secret storage",
	Run: func(cmd *cobra.Command, args []string) {
		v := secret.FileVault(encodingKey, secretsFile())
		key := args[0]
		err := v.Remove(key)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("removed key (%s)\n", key)
	},
}

func init() {
	RootCmd.AddCommand(removeCmd)
}
