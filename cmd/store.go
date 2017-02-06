// Store is the main command for handing gzr relevant data to inform it about
// infrastructure and artifacts that cannot be discovered automatically

package cmd

import (
	"fmt"
	"os"

	"github.com/bypasslane/gzr/comms"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var registeredInterfaces = map[string]comms.ImageStorageInterface{
	"etcd": &comms.EtcdImageStorer{},
}

var storeCmd = &cobra.Command{
	Use:   "store IMAGE_NAME METADATA_PATH",
	Short: "Store metadata about an image for gzr to track",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("Must provide IMAGE_NAME and METADATA_PATH")
			os.Exit(1)
		}
		var storer comms.ImageStorageInterface
		if registeredStore, ok := registeredInterfaces[viper.GetString("datastore.type")]; !ok {
			fmt.Printf("%s is not a valid datastore type", viper.GetString("datastore.type"))
			os.Exit(1)
		} else {
			storer = registeredStore
		}
		err := storer.StoreImage(args[0], args[1])
		if err != nil {
			fmt.Printf("Error storring image: %s", err.Error())
			os.Exit(1)
		}
	},
}

var getCmd = &cobra.Command{
	Use: "get IMAGE_NAME",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Must provide IMAGE_NAME")
			os.Exit(1)
		}
		var storer comms.ImageStorageInterface
		if registeredStore, ok := registeredInterfaces[viper.GetString("datastore.type")]; !ok {
			fmt.Printf("%s is not a valid datastore type", viper.GetString("datastore.type"))
			os.Exit(1)
		} else {
			storer = registeredStore
		}
		images, err := storer.GetImage(args[0])
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			os.Exit(1)
		}
		fmt.Printf("%+v\n", images)
	},
}

func init() {
	RootCmd.AddCommand(storeCmd)
	RootCmd.AddCommand(getCmd)
}
