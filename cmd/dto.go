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
	dtoFields   string
	dtoPackage  string
	dtoValidate bool
)

var dtoCmd = &cobra.Command{
	Use:   "dto [name]",
	Short: "Generate a DTO class",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		raw := args[0]
		if raw == "" {
			return errors.New("dto name is required")
		}
		name := exportName(raw)

		if err := requireSpringProject(); err != nil {
			return err
		}

		pkg := getEffectivePackage(".", installPackage, dtoPackage)

		dir := filepath.Join("src", "main", "java", filepath.Join(strings.Split(pkg, ".")...), "dto")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("création des dossiers: %w", err)
		}

		filePath := filepath.Join(dir, name+"Dto.java")
		content := generatePojoContent(name, pkg, "dto", "Dto", dtoFields, dtoValidate)
		if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
			return fmt.Errorf("écriture du fichier dto: %w", err)
		}

		Success("Created dto: %s\n", filePath)
		return nil
	},
}

func init() {
	dtoCmd.Flags().StringVar(&dtoFields, "fields", "", "fields like name:String,age:int")
	dtoCmd.Flags().StringVarP(&dtoPackage, "package", "p", "", "Override base package (ex: com.monentreprise.monprojet)")
	dtoCmd.Flags().BoolVar(&dtoValidate, "validate", false, "add Jakarta validation annotations")
	makeCmd.AddCommand(dtoCmd)
}
