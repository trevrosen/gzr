// Image is the main command for handing gzr relevant data to inform it about
// image artifacts that cannot be discovered automatically

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/bypasslane/gzr/comms"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var imageCmd = &cobra.Command{
	Use:   "image (store|get|delete)",
	Short: "manage information about images",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		storeType := viper.GetString("datastore.type")
		if storeType == "" {
			er("Must supply a datastore type in config file")
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
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		imageStore.Cleanup()
	},
}

var storeCmd = &cobra.Command{
	Use:   "store IMAGE_NAME:VERSION METADATA_PATH",
	Short: "Store metadata about an image for gzr to track",
	Long: `Used to store metadata about an image for gzr to track. The name must be formatted as NAME:VERSION.
Repeated store calls with the same VERSION on the same day will only keep the newest, but different days will be stored.
In short, only one version per day is allowed.

The structure of the JSON at the METADATA_PATH should be as follows:
{
    "git-commit": <string>,
    "git-tag": [<string>, <string>, ...],
    "git-annotation": [<string>, <string>, ...],
    "git-origin": <string>,
    "created-at": <ISO 8601-formatted string>
}`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			erBadUsage(fmt.Sprintf("Must provide IMAGE_NAME:VERSION and METADATA_PATH"), cmd)
		}
		reader, err := os.Open(args[1])
		if err != nil {
			er(fmt.Sprintf("Could not read metadata file"))
		}
		meta, err := comms.CreateMeta(reader)
		if err != nil {
			er(fmt.Sprintf("%s", err.Error()))
		}
		err = imageStore.Store(args[0], meta)
		if err != nil {
			er(fmt.Sprintf("Error storring image: %s", err.Error()))
		}
		fmt.Printf("Stored %s\n", args[0])
	},
}

var getCmd = &cobra.Command{
	Use:   "get IMAGE_NAME",
	Short: "Get data about the stored images under a particular name",
	Long: `Get all metadata about the stored images under a particular name,
including all versions held within gzr`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			erBadUsage(fmt.Sprintf("Must provide IMAGE_NAME"), cmd)
		}
		images, err := imageStore.List(args[0])
		if err != nil {
			er(fmt.Sprintf("Error: %s", err.Error()))
		}
		images.SerializeForCLI(os.Stdout)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete IMAGE_NAME:VERSION",
	Short: "Delete metadata about an image:version within gzr",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			erBadUsage(fmt.Sprintf("Must provide IMAGE_NAME:VERSION"), cmd)
		}
		splitName := strings.Split(args[0], ":")
		if len(splitName) != 2 {
			er(fmt.Sprintf("IMAGE_NAME must be formatted as NAME:VERSION and must contain only the seperating colon"))
		}
		err := imageStore.Delete(args[0])
		if err != nil {
			er(fmt.Sprintf("%s", err.Error()))
		}
		fmt.Printf("Deleted %s\n", args[0])
	},
}

func init() {
	imageCmd.AddCommand(storeCmd)
	imageCmd.AddCommand(getCmd)
	imageCmd.AddCommand(deleteCmd)
	RootCmd.AddCommand(imageCmd)
}
