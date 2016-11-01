package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the given service",
	Long:  `TODO: sample usage, multi-line explanation`,

	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 1:
			// - consume service config file
			// - check for Kubernetes availability
			// - start service in current context
		default:
			er("must specify service")
		}
		notify(fmt.Sprintf("starting '%s' in active context", args[0]))
	},
}

func init() {
	serviceCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
