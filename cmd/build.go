package cmd

import (
	"github.com/bypasslane/gzr/comms"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build [DOCKER ARGS...]",
	Short: "Wrapper around `docker build` to produce Docker artifacts as well as register data with gzr",
	Long: `Wrapper around "docker build" to produce Docker artifacts as well as register data with gzr
tagging (-t NAME:TAG) is required for metadata storage`,
	Run: func(cmd *cobra.Command, args []string) {
		// require tag option
		if imageTag == "" {
			er("Must provide --tag/-t flag with NAME:TAG")
		}
		// Add tag back to the docker args because it is pulled out since it is a flag
		dockerArgs := append(args, []string{"-t", imageTag}...)
		err := comms.BuildDocker(dockerArgs...)
		if err != nil {
			er(err.Error())
		}
		meta, err := comms.BuildMetadata()
		if err != nil {
			er(err.Error())
		}
		err = comms.PushDocker(imageTag)
		if err != nil {
			er(err.Error())
		}
		err = imageStore.Store(imageTag, meta)
		if err != nil {
			er(err.Error())
		}
	},
}

func init() {
	RootCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringVarP(&imageTag, "tag", "t", "", "Name and a tag in the NAME:TAG format")
}
