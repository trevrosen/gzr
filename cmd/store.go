package cmd

import (
	"github.com/bypasslane/gzr/comms"
	"github.com/spf13/cobra"
)

var storeCmd = &cobra.Command{
	Use:   "store [IMAGE_ID] [METADATA_PATH]",
	Short: "Store metadata about an image for gzr to track",
	Run: func(cmd *cobra.Command, args []string) {
		comms.StoreImage(args)
	},
}

func init() {
	RootCmd.AddCommand(storeCmd)
}
