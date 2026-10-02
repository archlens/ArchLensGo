package cmd

import (
	"maps"
	"path/filepath"
	"slices"

	"github.com/archlens/ArchLens/caching"
	"github.com/archlens/ArchLens/input"
	"github.com/archlens/ArchLens/mermaid"
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
		p, _ := cmd.Flags().GetBool("parse-full")
		markdown, _ := cmd.Flags().GetBool("markdown")

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

		absRoot, err := filepath.Abs(res.RootFolder)
		if err != nil {
			Sugar.Errorf("Error resolving root folder: %v", err)
			return
		}
		cachePath := filepath.Join(absRoot, ".archlens", "cache.json")
		cache := caching.LoadCache(cachePath)

		viewFiles := make(map[string][]string)
		union := make(map[string]bool)
		for name, view := range res.Views {
			files, err := input.GetFiles(&view, res.RootFolder)
			if err != nil {
				Sugar.Errorf("Error when trying to get files for %s: %v", name, err)
				continue
			}
			for i, f := range files {
				files[i] = filepath.ToSlash(filepath.Clean(f))
				union[files[i]] = true
			}
			viewFiles[name] = files
		}

		if err := cache.Refresh(absRoot, slices.Sorted(maps.Keys(union)), res.RunCommand, p); err != nil {
			Sugar.Errorf("Error refreshing cache: %v", err)
			return
		}

		// Graphs are cheap to build from cached results, so rebuild one per view.
		for name, files := range viewFiles {
			Sugar.Debugf("%s: %d files: %v", name, len(files), files)
			if len(files) == 0 {
				continue
			}
			g := cache.Graph(files)
			mermaid.Render(g, name, res.SaveLocation, markdown)
			Sugar.Infof("%s: %d nodes, %d edges", name, g.TotalNodes, g.TotalEdges)
		}

		if err := cache.Save(cachePath); err != nil {
			Sugar.Errorf("Error saving cache: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)
	renderCmd.Flags().BoolP("parse-full", "p", false, "A flag to fully parse any project again even if cache is present")
	renderCmd.Flags().Bool("markdown", false, "Export diagrams in markdown files")
}
