package cmd

import (
	"fmt"

	"github.com/bypasslane/gzr/comms"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var buildCmd = &cobra.Command{
	Use:   "build [DOCKER ARGS...]",
	Short: "Wrapper around `docker build` to produce Docker artifacts as well as register data with gzr",
	Long:  `Wrapper around "docker build" to produce Docker artifacts as well as register data with gzr`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// TODO: Dry this for setup between build/image
		storeType := viper.GetString("datastore.type")
		if storeType == "" {
			er("Must supply a datastore type in config file")
		}

		if viper.GetString("repository") == "" {
			er("Must provide \"repository\" setting in config file")
		}

		var storeCreator func() (comms.GzrMetadataStore, error)
		if creator, ok := registeredInterfaces[storeType]; !ok {
			er(fmt.Sprintf("%s is not a valid datastore type", storeType))
		} else {
			storeCreator = creator
		}
		newStore, err := storeCreator()
		if err != nil {
			er(fmt.Sprintf("%s", err.Error()))
		}
		imageStore = newStore

		buildEnv := viper.GetString("build_env")
		if buildEnv == "test" {
			imageManager = comms.NewDefaultMockManager()
		} else {
			imageManager = comms.NewDockerManager()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := buildHandler(args, imageManager)
		if err != nil {
			er(err.Error())
		}
	},
}

// buildHander handles the arguments from running a build command.
// The steps involved are as follows: Build image, create the metadata blob
// that accompanies the image, create the tag for docker, use a transaction
// to store the metadata and push the image
func buildHandler(args []string, manager comms.ImageManager) error {
	err := manager.Build(args...)
	if err != nil {
		return err
	}
	meta, err := comms.NewImageMetadata()
	if err != nil {
		return err
	}
	tag, err := comms.GetDockerTag()
	if err != nil {
		return err
	}
	err = imageStore.StartTransaction()
	if err != nil {
		return err
	}
	err = imageStore.Store(tag, meta)
	if err != nil {
		return err
	}
	err = manager.Push(tag)
	if err != nil {
		return err
	}
	err = imageStore.CommitTransaction()
	if err != nil {
		return err
	}
	return nil
}

func init() {
	RootCmd.AddCommand(buildCmd)
}
