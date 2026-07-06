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
	repositoryPackage string
	repositoryUUID    bool
)

var repositoryCmd = &cobra.Command{
	Use:   "repository [name]",
	Short: "Generate a JpaRepository interface for an entity",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		raw := args[0]
		if raw == "" {
			return errors.New("repository name is required")
		}
		name := exportName(raw)

		if err := requireSpringProject(); err != nil {
			return err
		}

		pkg := getEffectivePackage(".", installPackage, repositoryPackage)

		dir := filepath.Join("src", "main", "java", filepath.Join(strings.Split(pkg, ".")...), "repository")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("création des dossiers: %w", err)
		}

		filePath := filepath.Join(dir, name+"Repository.java")
		idType := "Long"
		if repositoryUUID {
			idType = "java.util.UUID"
		}
		
		content, err := renderTemplate("repository", struct {
			Pkg    string
			Entity string
			IdType string
		}{Pkg: pkg, Entity: name, IdType: idType})
		if err != nil {
			return fmt.Errorf("rendu du template repository: %w", err)
		}

		if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
			return fmt.Errorf("écriture du fichier repository: %w", err)
		}

		Success("Created repository: %s\n", filePath)
		return nil
	},
}

func init() {
	makeCmd.AddCommand(repositoryCmd)
	repositoryCmd.Flags().StringVarP(&repositoryPackage, "package", "p", "", "Override base package (ex: com.monentreprise.monprojet)")
	repositoryCmd.Flags().BoolVar(&repositoryUUID, "uuid", false, "use UUID instead of Long for the primary key")
}
