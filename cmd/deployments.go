package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/bypasslane/gzr/comms"
	"github.com/spf13/cobra"
)

// Package-global k8s connection
var k8sConn *comms.K8sConnection

// deploymentsCmd represents the deployments command
var deploymentsCmd = &cobra.Command{
	Use:   "deployments [get|list|update|] [args]",
	Short: "Manage k8s Deployments",
	Long: `Used to get information on single Deployments or all Deployments in a cluster

deployments list
deployments get <DEPLOYMENT NAME>
deployments update <DEPLOYMENT_NAME> <CONTAINER_NAME> <IMAGE>
	`,
	PreRun: func(cmd *cobra.Command, args []string) {
		var connErr error
		k8sConn, connErr = comms.NewK8sConnection(namespace)
		if connErr != nil {
			// TODO: figure out the Cobra way to handle this
			msg := fmt.Sprintf("problem establishing k8s connection: %s", connErr)
			er(msg)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			switch args[0] {
			case "list":
				displayListDeployments(namespace)
			case "get":
				if len(args) < 2 {
					er("'get' must be called with a Deployment name")
				}
				displayGetDeployment(namespace, args[1])
			case "update":
				updateDeployment(namespace, args[1], args[2], args[3])
			} // end switch
		} else {
			fmt.Println("Not enough arguments")
			fmt.Println(cmd.Use)
		}
	},
}

// updateDeployment updates a Deployment container with the info described by the DeploymentContainerInfo argument
func updateDeployment(namespace string, deploymentName string, containerName string, image string) {
	dci := comms.DeploymentContainerInfo{
		Namespace:      namespace,
		DeploymentName: deploymentName,
		ContainerName:  containerName,
		Image:          image,
	}
	deployment, err := k8sConn.UpdateDeployment(dci)

	if err != nil {
		msg := fmt.Sprintf("there was a problem updating container '%s' on deployment '%s' - %s", containerName, deploymentName, err)
		er(msg)
	}
	deployment.SerializeForCLI(os.Stdout)
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
	dlist, err := k8sConn.ListDeployments(namespace)
	if err != nil {
		log.Fatalln("Error retrieving Deployments: %s", err)
	}
	for _, deployment := range dlist.Deployments {
		deployment.SerializeForCLI(os.Stdout)
	}
}

func init() {
	RootCmd.AddCommand(deploymentsCmd)
	deploymentsCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "namespace to look for Deployments in")
}
