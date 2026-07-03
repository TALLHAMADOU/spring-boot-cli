package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var repositoryPackage string

var repositoryCmd = &cobra.Command{
	Use:   "repository [name]",
	Short: "Generate a JpaRepository interface for an entity",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		raw := args[0]
		if raw == "" {
			fmt.Fprintln(os.Stderr, "repository name is required")
			return
		}
		name := exportName(raw)

		if !isSpringProject(".") {
			fmt.Fprintln(os.Stderr, "Erreur: Lancez cette commande dans un projet Spring Boot (présence de pom.xml ou build.gradle)")
			os.Exit(1)
		}

		pkg := getEffectivePackage(".", installPackage, repositoryPackage)

		dir := filepath.Join("src", "main", "java", filepath.Join(strings.Split(pkg, ".")...), "repository")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "failed to create directories: %v\n", err)
			return
		}

		filePath := filepath.Join(dir, name+"Repository.java")
		content, err := renderTemplate("repository", struct {
			Pkg    string
			Entity string
		}{Pkg: pkg, Entity: name})
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to render repository template: %v\n", err)
			return
		}

		if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write repository file: %v\n", err)
			return
		}

		fmt.Printf("Created repository: %s\n", filePath)
	},
}

func init() {
	makeCmd.AddCommand(repositoryCmd)
	repositoryCmd.Flags().StringVarP(&repositoryPackage, "package", "p", "", "Override base package (ex: com.monentreprise.monprojet)")
}
