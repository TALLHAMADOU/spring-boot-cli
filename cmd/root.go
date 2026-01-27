/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "spring-cli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// Version is set at build time using -ldflags "-X 'github.com/your/module/cmd.Version=1.2.3'"
var Version = "dev"
var showVersion bool

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.spring-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// use Cobra's built-in Version handling for `--version`
	rootCmd.Version = Version

	// add short `-v` alias to show the CLI version (avoids conflicting with `version` subcommand)
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "cli-version", "v", false, "show spring-cli version")
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if showVersion {
			fmt.Printf("spring-cli %s\n", Version)
			os.Exit(0)
		}
	}
}
