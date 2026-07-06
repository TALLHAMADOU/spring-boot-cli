package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var exceptionPackage string

var exceptionCmd = &cobra.Command{
	Use:   "exception-handler",
	Short: "Generate a global @ControllerAdvice exception handler",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireSpringProject(); err != nil {
			return err
		}

		pkg := getEffectivePackage(".", installPackage, exceptionPackage)

		dir := filepath.Join("src", "main", "java", filepath.Join(strings.Split(pkg, ".")...), "exception")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("création des dossiers: %w", err)
		}

		filePath := filepath.Join(dir, "GlobalExceptionHandler.java")

		content, err := renderTemplate("exception_handler", struct{ Pkg string }{Pkg: pkg})
		if err != nil {
			return fmt.Errorf("rendu du template exception_handler: %w", err)
		}

		if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
			return fmt.Errorf("écriture du fichier exception_handler: %w", err)
		}

		Success("Created exception handler: %s\n", filePath)
		return nil
	},
}

func init() {
	exceptionCmd.Flags().StringVarP(&exceptionPackage, "package", "p", "", "Override base package")
	makeCmd.AddCommand(exceptionCmd)
}
