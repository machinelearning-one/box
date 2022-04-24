/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

// remoteCmd represents the remote command
var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Lists all the commands available through the remote",
	Long:  `Lists all the commands available through the remote.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Make a request to the remote
		resp, err := http.Get("http://resources.machinelearning.one/box/keys.json")
		if err != nil {
			fmt.Println("Error fetching remote")
			return
		}
		defer resp.Body.Close()
		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading remote response")
			return
		}
		var commands map[string]string
		err = json.Unmarshal(body, &commands)
		if err != nil {
			fmt.Println("Error parsing remote response")
			return
		}
		fmt.Println("Available commands:")
		for key, description := range commands {
			fmt.Println("-", key, ":", description)
		}
	},
}

func init() {
	lsCmd.AddCommand(remoteCmd)
}
