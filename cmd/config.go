package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/j178/leetgo/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show leetgo config dir",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.Get().ConfigDir())
	},
}

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt a sensitive string to be used in config file",
	RunE: func(cmd *cobra.Command, args []string) error {
		prompt := &survey.Password{
			Message: "Enter the string to be encrypted",
		}
		var input string
		err := survey.AskOne(prompt, &input)
		if err != nil {
			return err
		}
		encrypted, err := config.Encrypt(input)
		if err != nil {
			return err
		}
		fmt.Println()
		fmt.Println(encrypted)
		return nil
	},
}

func init() {
	configCmd.AddCommand(encryptCmd)
}
