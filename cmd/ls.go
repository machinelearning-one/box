/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"box/core"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Lists all the commands",
	Long:  `Lists all the commands in the box`,
	Run: func(cmd *cobra.Command, args []string) {
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
		fmt.Println("Available commands:")
		for key := range boxes {
			fmt.Println("-", key)
		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
