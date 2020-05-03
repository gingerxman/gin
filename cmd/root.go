package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var logoWithDesc = `
   _____ _
  / ____(_)
 | |  __ _ _ __
 | | |_ | | '_ \
 | |__| | | | | |
  \_____|_|_| |_|  ginger micro-service cli v%s

`

var logo = `
   _____ _
  / ____(_)
 | |  __ _ _ __
 | | |_ | | '_ \
 | |__| | | | | |
  \_____|_|_| |_| ginger micro-service cli v%s

`

var rootCmd = &cobra.Command{
	Use:   "gin",
	Short: "Hugo is a very fast static site generator",
	Long: fmt.Sprintf(logoWithDesc, VERSION),
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func PrintBanner() {
	fmt.Println(fmt.Sprintf(logo, VERSION))
}