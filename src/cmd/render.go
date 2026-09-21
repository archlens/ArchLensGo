package cmd

import (
	"github.com/archlens/ArchLens/input"
	"github.com/archlens/ArchLens/parsers"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	ts_python "github.com/tree-sitter/tree-sitter-python/bindings/go"
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
		logger, _ := zap.NewProduction()
		defer logger.Sync()
		sugar := logger.Sugar()

		var configPath string
		if len(args) == 0 {
			configPath = "archlens.json"
		} else {
			configPath = args[0]
		}

		res, err := input.Load(configPath)
		if err != nil {
			sugar.Errorf("Error when trying to load configuration: %v", err)
			return
		}
		sugar.Infof("Config: %+v", *res)

		viewFiles := make(map[string][]string)
		for name, view := range res.Views {
			files, err := input.GetFiles(&view, res.RootFolder)
			if err != nil {
				sugar.Errorf("Error when trying to get files for %s: %v", name, err)
			}
			viewFiles[name] = files
		}

		var data [][]byte
		var errs []error
		for name, files := range viewFiles {
			sugar.Debugf("%s: %d files: %v", name, len(files), files)
			if len(files) == 0 {
				continue
			}

			data, errs = input.ReadFiles(files)
			if len(data) == 0 || data[0] == nil {
				continue
			}

			for _, err := range errs {
				if err != nil {
					sugar.Error(err)
				}
			}
			parsers.Parse(ts_python.Language(), data[0])
		}

	},
}

func init() {
	rootCmd.AddCommand(renderCmd)
}
