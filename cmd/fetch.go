package cmd

import (
	"box/core"
	"box/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetches a command from the remote",
	Long:  `Fetches a command from the remote.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get the first argument as the key
		key := args[0]
		// Make a request to the remote
		resp, err := http.Get("http://resources.machinelearning.one/box/commands/" + key + ".json")
		if err != nil {
			fmt.Println("Error fetching remote, please make sure you have entered a valid key")
			return
		}
		defer resp.Body.Close()
		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading remote response")
			return
		}
		var remoteConfig core.BoxConfig
		err = json.Unmarshal(body, &remoteConfig)
		if err != nil {
			fmt.Println("Error parsing remote response")
			return
		}
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
		alt, _ := cmd.Flags().GetString("alt-key")
		if alt != "" {
			key = alt
		}
		// Check if the key already exists
		if _, ok := boxes[key]; ok {
			fmt.Println("Key already exists, please choose another")
			return
		}
		// Check if image is already present else pull it
		err = utils.TouchImage(remoteConfig.Image)
		if err != nil {
			fmt.Println("Error fetching image")
			return
		}
		// Add box config to map
		boxes[key] = remoteConfig
		// Marshal map to json
		json, err := json.Marshal(boxes)
		if err != nil {
			fmt.Println("Error converting fetched information to valid configuration")
			return
		}
		// Write json to file
		err = ioutil.WriteFile(".box", json, 0644)
		if err != nil {
			fmt.Println("Error writing command to .box file")
			return
		}
		fmt.Println("Remote command successfully added")
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
	// Add optional alt-key flag
	fetchCmd.Flags().StringP("alt-key", "k", "", "Alternative key to use to avoid conflicts")
}
