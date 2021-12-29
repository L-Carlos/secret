package cobra

import (
	"fmt"

	"github.com/L-Carlos/secret"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "sets a secret in your secret storage",
	Run: func(cmd *cobra.Command, args []string) {
		v := secret.FileVault(encodingKey, secretsFile())
		if len(args) != 2 {
			fmt.Println("wrong number of arguments, expected [key] [value]")
			return
		}
		key, value := args[0], args[1]
		err := v.Set(key, value)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Value Set for key (%s)\n", key)
	},
}

func init() {
	RootCmd.AddCommand(setCmd)
}
