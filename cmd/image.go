// Image is the main command for handing gzr relevant data to inform it about
// image artifacts that cannot be discovered automatically

package cmd

import (
	"fmt"
	"os"

	"github.com/bypasslane/gzr/comms"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var registeredInterfaces = map[string]comms.GozerMetadataStore{
	"etcd": &comms.EtcdStorage{},
}

var imageCmd = &cobra.Command{
	Use:   "image (store|get)",
	Short: "Base command for image data manipulation",
}

var storeCmd = &cobra.Command{
	Use:   "store IMAGE_NAME:VERSION METADATA_PATH",
	Short: "Store metadata about an image for gzr to track",
	Long: `Used to store metadata about an image for gzr to track. The name must be formatted as NAME:VERSION.
Repeated store calls with the same VERSION on the same day will only keep the newest, but different days will be stored.
In short, only one version per day is allowed.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("Must provide IMAGE_NAME:VERSION and METADATA_PATH")
			os.Exit(1)
		}
		var store comms.GozerMetadataStore
		if registeredStore, ok := registeredInterfaces[viper.GetString("datastore.type")]; !ok {
			fmt.Printf("%s is not a valid datastore type", viper.GetString("datastore.type"))
			os.Exit(1)
		} else {
			store = registeredStore
		}
		reader, err := os.Open(args[1])
		if err != nil {
			fmt.Printf("Could not read metadata file")
			os.Exit(1)
		}
		meta, err := comms.CreateMeta(reader)
		if err != nil {
			fmt.Printf("%s", err.Error())
			os.Exit(1)
		}
		err = store.Store(args[0], meta)
		if err != nil {
			fmt.Printf("Error storring image: %s", err.Error())
			os.Exit(1)
		}
	},
}

var getCmd = &cobra.Command{
	Use:   "get IMAGE_NAME",
	Short: "Get data about the stored images under a particular name",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Must provide IMAGE_NAME")
			os.Exit(1)
		}
		var store comms.GozerMetadataStore
		if registeredStore, ok := registeredInterfaces[viper.GetString("datastore.type")]; !ok {
			fmt.Printf("%s is not a valid datastore type", viper.GetString("datastore.type"))
			os.Exit(1)
		} else {
			store = registeredStore
		}
		images, err := store.List(args[0])
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			os.Exit(1)
		}
		fmt.Printf("%+v\n", images)
	},
}

func init() {
	imageCmd.AddCommand(storeCmd)
	imageCmd.AddCommand(getCmd)
	RootCmd.AddCommand(imageCmd)
}
