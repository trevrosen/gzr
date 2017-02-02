// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/bypasslane/gzr/comms"
	"github.com/spf13/cobra"
)

// Global connection used by this command
var k8sConn *comms.K8sConnection

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [k8s RESOURCE TYPE]",
	Short: "retrieve information about Kubernetes resources",
	Long: `Works similarly to kubectl get but with more opinionated output:

gzr get deployments`,
	PreRun: func(cmd *cobra.Command, args []string) {
		var connErr error
		k8sConn, connErr = comms.NewK8sConnection()
		if connErr != nil {
			// TODO: figure out the Cobra way to handle this
			fmt.Println("Error establishing k8s connection: ", connErr)
			os.Exit(1)
		}
	},

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			// TODO: make this use Cobra usage message
			fmt.Println("Resource type can't be blank")
			os.Exit(1)
		}
		switch args[0] {
		case "deployments":
			// TODO: allow getting a single Deployment by name
			// TODO: take in namespace from flag
			activeDeployments, err := k8sConn.Deployments("default")
			if err != nil {
				log.Fatalln("Error retrieving Deployments: %s", err)
			}
			for _, deployment := range activeDeployments {
				deployment.SerializeForCLI(os.Stdout)
			}
		} // end switch

	},
}

func init() {
	RootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
