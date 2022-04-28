package cmd

import (
	"box/core"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// aliasCmd represents the alias command
var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Creates alias for existing commands",
	Long:  `Creates alias for existing commands.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		rest := args[1:]
		alias, _ := cmd.Flags().GetString("alias")
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
		// Since the key exists, we can create an alias
		// Get home directory
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Could not get home directory, make sure you have sufficient permissions")
			return
		}
		// Add the rest of the arguments to the key separated by spaces
		if len(rest) > 0 {
			key = key + " " + strings.Join(rest, " ")
		}
		// Create the alias line
		aliasLine := fmt.Sprintf("alias %s='%s %s'\n", alias, "box run", key)
		// Check which shell config files exist
		shells := []string{"bash", "zsh", "fish"}
		configExists := false
		for _, shell := range shells {
			file := fmt.Sprintf("%s/.%src", home, shell)
			if _, err := os.Stat(file); err == nil {
				configExists = true
				// Add alias to shell config file
				fmt.Printf("Adding alias to %s\n", shell)
				// Read the file
				bytes, err := ioutil.ReadFile(file)
				if err != nil {
					fmt.Println("Error reading", file)
					return
				}
				// Add the alias line to the file
				newBytes := append(bytes, []byte(aliasLine)...)
				// Write the new file
				err = ioutil.WriteFile(file, newBytes, 0644)
				if err != nil {
					fmt.Println("Error writing", file)
					return
				}
			}
		}
		if !configExists {
			fmt.Println("Could not find any shell config files")
			return
		}
		fmt.Println("Alias created")
	},
}

func init() {
	rootCmd.AddCommand(aliasCmd)
	// Get the alias and mark it as required
	aliasCmd.Flags().StringP("alias", "a", "", "The alias to create")
	aliasCmd.MarkFlagRequired("alias")
}
