package cmd

import (
	"fmt"
	
	"github.com/spf13/cobra"
)

const VERSION = "0.1"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Gin",
	
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(fmt.Sprintf("Gin v%s", VERSION))
	},
}


func init() {
	rootCmd.AddCommand(versionCmd)
}