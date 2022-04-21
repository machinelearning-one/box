package cmd

import (
	"box/core"
	"box/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Removes a command from the box",
	Long:  `Removes an existing command from the box, removes the image if not used by any other command`,
	Run: func(cmd *cobra.Command, args []string) {
		key, _ := cmd.Flags().GetString("key")
		keepImage, _ := cmd.Flags().GetBool("keep-image")
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
		if !keepImage {
			// Check if image is used by any other command
			used := false
			for id, box := range boxes {
				// If box is not the one to be removed and the image is the same
				if box.Image == boxes[key].Image && id != key {
					used = true
					break
				}
			}
			if !used {
				// Remove image
				err = utils.RemoveImage(boxes[key].Image)
				if err != nil {
					fmt.Println("Error removing image")
					return
				}
			}
		}
		// Remove the key
		delete(boxes, key)
		// Write the updated boxes to the .box file
		bytes, err = json.Marshal(boxes)
		if err != nil {
			fmt.Println("Error converting to valid configuration")
			return
		}
		err = ioutil.WriteFile(".box", bytes, 0644)
		if err != nil {
			fmt.Println("Error updating .box file")
			return
		}
		fmt.Println("Command successfully removed")
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
	// Get the key and mark it as required
	rmCmd.Flags().StringP("key", "k", "", "Key of the command to remove")
	rmCmd.MarkFlagRequired("key")
	// Get optional keep-image flag
	rmCmd.Flags().BoolP("keep-image", "i", false, "Keep the image after removing the command")
}
