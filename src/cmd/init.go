package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/archlens/ArchLens/input"
	"github.com/archlens/ArchLens/parsers"
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
	Use:   "init [language]",
	Short: "Initializes the project with a base archlens.json file and optionally [language] based ast_parser",
	Long: `Initializes the project with a base archlens.json file as well as an embedded ast_parser in the specified optional [language] argument.
	The project currently ships with python and go ast_parsers. Users can create their own parsers by using the specification that can be found in the docs folder
	on the github page. Archlens communicates with the parser using a unified json struct through IPC so it's easy to build your own. 
	
	The location for the archlens.json file as well as the ast_parsers can be specified using the -l flag the parser folder is always initialized next to
	the archlens.json file.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		l, _ := cmd.Flags().GetString("location")

		var language string
		if len(args) > 0 {
			language = args[0]
		}

		path := "./archlens.json"
		if l != "" {
			path = l
		}

		info, err := os.Stat(path)
		switch {
		case err == nil && info.IsDir():
			// Existing directory: put archlens.json inside it.
			path = filepath.Join(path, "archlens.json")
		case filepath.Ext(path) == "":
			// No extension: treat it as a directory path.
			path = filepath.Join(path, "archlens.json")
		case !strings.EqualFold(filepath.Ext(path), ".json"):
			Sugar.Errorf("config path %q must be a .json file or a directory", path)
			return
		}

		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			Sugar.Errorf("Could not create %s: %s", dir, err)
			return
		}

		data, err := json.Marshal(template)
		if err != nil {
			panic(err)
		}

		err = os.WriteFile(path, data, 0644)
		if err != nil {
			Sugar.Errorf("Could not write archlens.json: %s", err.Error())
			return
		}
		fmt.Println("Wrote archlens.json to current location")

		if language != "" {
		pickLanguage(language, filepath.Join(dir, "ast_parser"))
		}
	},
}

func pickLanguage(language, location string) {
	language = strings.ToLower(language)
	switch language {
	case "python":
		err := parsers.BundlePython(location)
		if err != nil {
			Sugar.Errorf("Problem bundling the ast_parser at location: %s\nwith error: %s",location ,err)
		}
	case "go":
		err := parsers.BundleGo(location)
		if err != nil {
			Sugar.Errorf("Problem bundling the ast_parser at location: %s\nwith error: %s",location ,err)
		}
	default: 
		Sugar.Errorf("Language is not supported out of the box yet please create your own AST-parser using the specification on the github repository")
	}
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("location", "l", "", "Choose the location to put the default archlens.json file (defaults to ./archlens.json)")
}
