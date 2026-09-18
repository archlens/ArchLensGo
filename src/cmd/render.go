package cmd

import (
	"fmt"

	"github.com/archlens/ArchLens/input"
	"github.com/spf13/cobra"
)

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var configPath string
		if len(args) == 0 {
			configPath = "archlens.json"
		} else {
			configPath = args[0]
		}
		
		res, err := input.Load(configPath)
		if err != nil {
			fmt.Printf("Error when trying to load configuration: %v\n", err)
			return
		}
		fmt.Printf("Config: %+v\n", *res)
		viewFiles := make(map[string][]string)
		for name, view := range res.Views {
			files, err := input.GetFiles(&view, res.RootFolder)
			if err != nil {
				fmt.Printf("Error when trying to get files for %s: %v\n", name, err)
			}
			viewFiles[name] = files
		}
		for name, files := range viewFiles {
			fmt.Printf("%s:\t%+s\n", name, files)
		}
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// renderCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// renderCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
