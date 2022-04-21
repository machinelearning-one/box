package cmd

import (
	"box/core"
	"box/exec"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs a command in a container",
	Long:  `Runs a command in a container.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get the first argument as the key
		key := args[0]
		rest := args[1:]
		var boxes map[string]core.BoxConfig
		bytes, err := ioutil.ReadFile(".box")
		if err != nil {
			fmt.Println("Error reading .box file")
			return
		}
		err = json.Unmarshal(bytes, &boxes)
		if err != nil {
			fmt.Println("The .box file contains invalid configuration")
			return
		}
		// Check if the key already exists
		if _, ok := boxes[key]; !ok {
			fmt.Println("Key does not exist, please choose another")
			return
		}
		// Get the box config
		config := boxes[key]
		// Convert the config to a box
		box, err := core.ConfigToBox(config, rest)
		if err != nil {
			fmt.Println("Error converting config to an acceptable container")
			return
		}
		// Run the box
		ctx := context.Background()
		client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			fmt.Println("Error creating a docker client")
			return
		}
		err = exec.Ephemeral(ctx, client, box)
		if err != nil {
			fmt.Println("Error running the container")
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
