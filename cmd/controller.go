package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var controllerPackage string

var controllerCmd = &cobra.Command{
	Use:   "controller [name]",
	Short: "Generate a REST controller",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		raw := args[0]
		if raw == "" {
			return errors.New("controller name is required")
		}
		name := exportName(raw)

		if err := requireSpringProject(); err != nil {
			return err
		}

		pkg := getEffectivePackage(".", installPackage, controllerPackage)

		dir := filepath.Join("src", "main", "java", filepath.Join(strings.Split(pkg, ".")...), "controller")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("création des dossiers: %w", err)
		}

		filePath := filepath.Join(dir, name+"Controller.java")

		crud, _ := cmd.Flags().GetBool("crud")
		entityName, _ := cmd.Flags().GetString("entity")
		validate, _ := cmd.Flags().GetBool("validate")

		var content string
		if crud {
			if entityName == "" {
				return errors.New("--entity is required for --crud")
			}
			e := exportName(entityName)

			// ensure repository and service exist
			if err := ensureRepository(pkg, e); err != nil {
				return fmt.Errorf("génération du repository: %w", err)
			}
			if err := ensureService(pkg, e+"Service", e); err != nil {
				return fmt.Errorf("génération du service: %w", err)
			}

			var renderErr error
			content, renderErr = renderTemplate("controller_crud", struct {
				Pkg         string
				Name        string
				Entity      string
				EntityLower string
				Validate    bool
			}{Pkg: pkg, Name: e, Entity: e, EntityLower: strings.ToLower(e), Validate: validate})
			if renderErr != nil {
				return fmt.Errorf("rendu du template controller: %w", renderErr)
			}
		} else {
			var renderErr error
			content, renderErr = renderTemplate("controller_basic", struct {
				Pkg       string
				Name      string
				NameLower string
			}{Pkg: pkg, Name: name, NameLower: strings.ToLower(name)})
			if renderErr != nil {
				return fmt.Errorf("rendu du template controller: %w", renderErr)
			}
		}

		if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
			return fmt.Errorf("écriture du fichier controller: %w", err)
		}

		Success("Created controller: %s\n", filePath)
		return nil
	},
}

func init() {
	controllerCmd.Flags().Bool("crud", false, "generate CRUD endpoints")
	controllerCmd.Flags().String("entity", "", "entity name to use for CRUD (required with --crud)")
	controllerCmd.Flags().StringVarP(&controllerPackage, "package", "p", "", "Override base package (ex: com.monentreprise.monprojet)")
	controllerCmd.Flags().Bool("validate", false, "add @Valid to request bodies")
	makeCmd.AddCommand(controllerCmd)
}
