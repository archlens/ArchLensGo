package cmd

import (
	"fmt"

	"github.com/archlens/ArchLens/input"
	"github.com/archlens/ArchLens/parsers"
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
		// probably move this and make it a bit more generalized

		var configPath string
		if len(args) == 0 {
			configPath = "archlens.json"
		} else {
			configPath = args[0]
		}

		res, err := input.Load(configPath)
		if err != nil {
			Sugar.Errorf("Error when trying to load configuration: %v", err)
			return
		}
		Sugar.Infof("Config: %+v", *res)

		viewFiles := make(map[string][]string)
		for name, view := range res.Views {
			files, err := input.GetFiles(&view, res.RootFolder)
			if err != nil {
				Sugar.Errorf("Error when trying to get files for %s: %v", name, err)
			}
			viewFiles[name] = files
		}

		// TODO: We could potentially even wait group the views 
		for name, files := range viewFiles {
			Sugar.Debugf("%s: %d files: %v", name, len(files), files)
			if len(files) == 0 {
				continue
			}

			results := parsers.GetASTs(files, res.RootFolder)

			for _, r := range results {
				if r.Err != nil {
					fmt.Printf("%s: error: %v\n", r.File, r.Err)
					continue
				}
				fmt.Printf("%s: got AST: %+v\n", r.File, r.AST)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)
}
