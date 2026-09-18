package cmd

import (

	"github.com/spf13/cobra"
)

// renderDiffCmd represents the renderDiff command
var renderDiffCmd = &cobra.Command{
	Use:   "renderDiff",
	Short: "Renders difference views comparing current branch to the base branch",
	Long: `Renders difference views comparing current branch to the base branch`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(renderDiffCmd)
}
