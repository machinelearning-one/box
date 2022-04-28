package cmd

import (
	"box/core"
	"box/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a command to the box",
	Long:  `Validates and adds a new command to the box, fetches image if not already present`,
	Run: func(cmd *cobra.Command, args []string) {
		err := utils.TouchBoxFile()
		if err != nil {
			fmt.Println("Error accessing or creating .box file")
			return
		}
		key, _ := cmd.Flags().GetString("key")
		image, _ := cmd.Flags().GetString("image")
		command, _ := cmd.Flags().GetString("command")
		no_mount, _ := cmd.Flags().GetBool("no-mount")
		gpu, _ := cmd.Flags().GetBool("gpu")

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
		if _, ok := boxes[key]; ok {
			fmt.Println("Key already exists, please choose another")
			return
		}
		// Check if image is already present else pull it
		err = utils.TouchImage(image)
		if err != nil {
			fmt.Println("Error fetching image")
			return
		}
		// Create new box config from flags
		box := core.BoxConfig{
			Image: image,
			Cmd:   []string{command},
			GPUs:  gpu,
		}
		if !no_mount {
			box.Target = "/workspace"
			box.WorkDir = "/workspace"
		}
		// Add box config to map
		boxes[key] = box
		// Marshal map to json
		json, err := json.Marshal(boxes)
		if err != nil {
			fmt.Println("Error converting provided flags to valid configuration")
			return
		}
		// Write json to file
		err = ioutil.WriteFile(".box", json, 0644)
		if err != nil {
			fmt.Println("Error writing command to .box file")
			return
		}
		fmt.Println("Command successfully added")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	// Get the key and mark it as required
	addCmd.Flags().StringP("key", "k", "", "The key of the command")
	addCmd.MarkFlagRequired("key")
	// Get the image name and mark it as required
	addCmd.Flags().StringP("image", "i", "", "The name of the image to use")
	addCmd.MarkFlagRequired("image")
	// Get the command name and mark it as required
	addCmd.Flags().StringP("command", "c", "", "The command to run")
	addCmd.MarkFlagRequired("command")
	// Add optional no mount flag
	addCmd.Flags().BoolP("no-mount", "m", false, "Do not map current directory to /workspace")
	// Add boolean flag whether to use gpus
	addCmd.Flags().BoolP("gpu", "g", false, "Use gpus")
}
