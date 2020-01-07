package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "gret",
	Short: "Resource and Environment testing tool implemented with Go",
	Long: `Resource and Environment testing tool implemented with Go.

Usualy usecase are file exist test and environment value exist test.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("root called")
	},
}
