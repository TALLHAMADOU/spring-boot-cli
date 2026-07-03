package cmd

import (
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
	Run: func(cmd *cobra.Command, args []string) {
		raw := args[0]
		if raw == "" {
			fmt.Fprintln(os.Stderr, "controller name is required")
			return
		}
		name := exportName(raw)

		if !isSpringProject(".") {
			fmt.Fprintln(os.Stderr, "Erreur: Lancez cette commande dans un projet Spring Boot (présence de pom.xml ou build.gradle)")
			os.Exit(1)
		}

		pkg := getEffectivePackage(".", installPackage, controllerPackage)

		dir := filepath.Join("src", "main", "java", filepath.Join(strings.Split(pkg, ".")...), "controller")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "failed to create directories: %v\n", err)
			return
		}

		filePath := filepath.Join(dir, name+"Controller.java")

		crud, _ := cmd.Flags().GetBool("crud")
		entityName, _ := cmd.Flags().GetString("entity")

		var content string
		if crud {
			if entityName == "" {
				fmt.Fprintln(os.Stderr, "--entity is required for --crud")
				return
			}
			e := exportName(entityName)

			// ensure repository and service exist
			if err := ensureRepository(pkg, e); err != nil {
				fmt.Fprintf(os.Stderr, "failed to ensure repository: %v\n", err)
			}
			if err := ensureService(pkg, e+"Service", e); err != nil {
				fmt.Fprintf(os.Stderr, "failed to ensure service: %v\n", err)
			}

			var renderErr error
			content, renderErr = renderTemplate("controller_crud", struct {
				Pkg         string
				Name        string
				Entity      string
				EntityLower string
			}{Pkg: pkg, Name: e, Entity: e, EntityLower: strings.ToLower(e)})
			if renderErr != nil {
				fmt.Fprintf(os.Stderr, "failed to render controller template: %v\n", renderErr)
				return
			}
		} else {
			var renderErr error
			content, renderErr = renderTemplate("controller_basic", struct {
				Pkg       string
				Name      string
				NameLower string
			}{Pkg: pkg, Name: name, NameLower: strings.ToLower(name)})
			if renderErr != nil {
				fmt.Fprintf(os.Stderr, "failed to render controller template: %v\n", renderErr)
				return
			}
		}

		if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write controller file: %v\n", err)
			return
		}

		fmt.Printf("Created controller: %s\n", filePath)
	},
}

func init() {
	controllerCmd.Flags().Bool("crud", false, "generate CRUD endpoints")
	controllerCmd.Flags().String("entity", "", "entity name to use for CRUD (required with --crud)")
	controllerCmd.Flags().StringVarP(&controllerPackage, "package", "p", "", "Override base package (ex: com.monentreprise.monprojet)")
	makeCmd.AddCommand(controllerCmd)
}
