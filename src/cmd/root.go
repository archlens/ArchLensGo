package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "archlens",
	Short: "ArchLens enables you to create architectural views of your codebase",
	Long:  `ArchLens enables you to create architectural views of your codebase`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Global logger for ease of use
var Sugar *zap.SugaredLogger

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	logger, _ := zap.NewProduction()
	defer func() {
		_ = logger.Sync()
	}()
	Sugar = logger.Sugar()

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
