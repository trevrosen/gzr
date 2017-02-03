package cmd

import (
	"fmt"
	"os"

	"github.com/bypasslane/gzr/comms"
	"github.com/spf13/cobra"
)

var storeCmd = &cobra.Command{
	Use:   "store IMAGE_ID METADATA_PATH",
	Short: "Store metadata about an image for gzr to track",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			// TODO: make this use Cobra usage message
			fmt.Println("Must provide IMAGE_ID and METADATA_PATH")
			os.Exit(1)
		}
		comms.StoreImage(args)
	},
}

func init() {
	RootCmd.AddCommand(storeCmd)
}
