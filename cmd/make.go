package cmd

import (
	"github.com/spf13/cobra"
)

var makeCmd = &cobra.Command{
	Use:   "make",
	Short: "Generate project components (entities, etc.)",
}

func init() {
	rootCmd.AddCommand(makeCmd)
}
