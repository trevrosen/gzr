package cmd

import (
	"fmt"

	"github.com/bypasslane/gzr/comms"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build [DOCKER ARGS...]",
	Short: "Wrapper around `docker build` to produce Docker artifacts as well as register data with gzr",
	Long:  `Add me`,
	Run: func(cmd *cobra.Command, args []string) {
		dockerArgs := append([]string{"build"}, args...)
		comms.BuildDocker(dockerArgs...)
		meta, err := comms.BuildMetadata()
		if err != nil {
			er(err.Error())
		}
		fmt.Println(meta)
	},
}

func init() {
	RootCmd.AddCommand(buildCmd)
}
