package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/archlens/ArchLens/input"
	"github.com/spf13/cobra"
)


var template input.Input = input.Input {
	Name: "ArchLens",	
	RootFolder: ".",
	Github: input.Github{
		Url: "https://github.com/archlens/ArchLens",
		Branch: "master",
	},	
	SaveLocation: "./diagrams/",
	Views: map[string]input.View {
		"completeView": input.View{
			Include: []string{"**/*.*","*.*"},
			Exclude: []string{"go.sum","go.mod"},
		},
	},
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		l, _ := cmd.Flags().GetString("location")

		data, err := json.Marshal(template)
		if err != nil {
			panic(err)
		}
	
		var path = "./archlens.json"
		if l != "" {
			path = l
		}

		err = os.WriteFile(path, data, 0644)
		if err != nil {
			Sugar.Errorf("Could not write archlens.json: %s", err.Error())
		}
		fmt.Println("Wrote archlens.json to current location")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("location", "l", "", "Choose the location to put the default archlens.json file (defaults to ./archlens.json)")
}
