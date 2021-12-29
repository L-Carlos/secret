package cobra

import (
	"fmt"

	"github.com/L-Carlos/secret"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "lists all secrets from your secret storage",
	Run: func(cmd *cobra.Command, args []string) {
		v := secret.FileVault(encodingKey, secretsFile())
		err := v.List()
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
