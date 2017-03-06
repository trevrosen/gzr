package cmd

import (
	"github.com/bypasslane/gzr/comms"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build [DOCKER ARGS...]",
	Short: "Wrapper around `docker build` to produce Docker artifacts as well as register data with gzr",
	Long:  `Wrapper around "docker build" to produce Docker artifacts as well as register data with gzr`,
	Run: func(cmd *cobra.Command, args []string) {
		builder := comms.NewDockerBuilder()
		err := buildImage(args, builder)
		if err != nil {
			er(err.Error())
		}
	},
}

func buildImage(args []string, builder comms.ImageBuilder) error {
	err := builder.Build(args...)
	if err != nil {
		return err
	}
	meta, err := comms.BuildMetadata()
	if err != nil {
		return err
	}
	tag, err := comms.GetDockerTag()
	if err != nil {
		return err
	}
	err = builder.Push(tag)
	if err != nil {
		return err
	}
	err = imageStore.Store(tag, meta)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	RootCmd.AddCommand(buildCmd)
}
