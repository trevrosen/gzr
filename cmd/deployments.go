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

// deploymentsCmd represents the deployments command
var deploymentsCmd = &cobra.Command{
	Use:   "deployments [get|list|update|] [args]",
	Short: "Manage k8s Deployments",
	Long: `Used to get information on single Deployments or all Deployments in a cluster

deployments list
deployments get <DEPLOYMENT NAME>
deployments update <DEPLOYMENT NAME> <IMAGE>
	`,
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
		if len(args) > 0 {
			switch args[0] {
			case "list":
				// TODO: take in namespace from flag
				displayListDeployments("default")
			case "get":
				if len(args) < 2 {
					er("'get' must be called with a Deployment name")
				}
				displayGetDeployment("default", args[1])
			} // end switch
		} // TODO: have an else here that shows usage
	},
}

// displayGetDeployment fetches
func displayGetDeployment(namespace string, deploymentName string) {
	deployment, err := k8sConn.GetDeployment(namespace, deploymentName)
	if err != nil {
		msg := fmt.Sprintf("there was a problem retrieving deployment '%s'", deploymentName)
		er(msg)
	}
	deployment.SerializeForCLI(os.Stdout)
}

// displayListDeployments fetches Deployments and prints them to the CLI
func displayListDeployments(namespace string) {
	activeDeployments, err := k8sConn.ListDeployments(namespace)
	if err != nil {
		log.Fatalln("Error retrieving Deployments: %s", err)
	}
	for _, deployment := range activeDeployments {
		deployment.SerializeForCLI(os.Stdout)
	}
}

func init() {
	RootCmd.AddCommand(deploymentsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deploymentsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deploymentsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
