package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "verifyenv",
	Short: "file and Environment verify tool implemented with Go",
	Long: `file and Environment verify tool implemented with Go.

Usualy usecase are file exist check and environment value exist check.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("root called")
	},
}
