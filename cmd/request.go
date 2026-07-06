package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	requestFields   string
	requestPackage  string
	requestValidate bool
)

var requestCmd = &cobra.Command{
	Use:   "request [name]",
	Short: "Generate a Request class (for controller inputs)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		raw := args[0]
		if raw == "" {
			return errors.New("request name is required")
		}
		name := exportName(raw)

		if err := requireSpringProject(); err != nil {
			return err
		}

		pkg := getEffectivePackage(".", installPackage, requestPackage)

		dir := filepath.Join("src", "main", "java", filepath.Join(strings.Split(pkg, ".")...), "request")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("création des dossiers: %w", err)
		}

		filePath := filepath.Join(dir, name+"Request.java")
		content := generatePojoContent(name, pkg, "request", "Request", requestFields, requestValidate)
		if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
			return fmt.Errorf("écriture du fichier request: %w", err)
		}

		Success("Created request: %s\n", filePath)
		return nil
	},
}

func init() {
	requestCmd.Flags().StringVar(&requestFields, "fields", "", "fields like name:String,age:int")
	requestCmd.Flags().StringVarP(&requestPackage, "package", "p", "", "Override base package (ex: com.monentreprise.monprojet)")
	requestCmd.Flags().BoolVar(&requestValidate, "validate", false, "add Jakarta validation annotations")
	makeCmd.AddCommand(requestCmd)
}
