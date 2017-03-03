package cmd

import (
	"fmt"

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
		builder := comms.NewDockerBuilder()
		args = append(args, []string{"-t", imageTag}...)
		err := buildImage(args, builder)
		if err != nil {
			er(err.Error())
		}
	},
}

func buildImage(args []string, builder comms.ImageBuilder) error {
	// Add tag back to the docker args because it is pulled out since it is a flag
	err := builder.Build(args...)
	if err != nil {
		return err
	}
	meta, err := comms.BuildMetadata()
	if err != nil {
		return err
	}
	err = builder.Push(imageTag)
	if err != nil {
		return err
	}
	fmt.Printf("imageTag: %s\n", imageTag)
	fmt.Printf("meta: %+v\n", meta)
	err = imageStore.Store(imageTag, meta)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	RootCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringVarP(&imageTag, "tag", "t", "", "Name and a tag in the NAME:TAG format")
}
