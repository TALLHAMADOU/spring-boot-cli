package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var dtoFields string
var dtoPackage string

var dtoCmd = &cobra.Command{
	Use:   "dto [name]",
	Short: "Generate a DTO class",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		raw := args[0]
		if raw == "" {
			Error("dto name is required")
			return
		}
		name := exportName(raw)

		if !isSpringProject(".") {
			Error("Erreur: Lancez cette commande dans un projet Spring Boot (présence de pom.xml ou build.gradle)")
			os.Exit(1)
		}

		pkg := getEffectivePackage(".", installPackage, dtoPackage)

		dir := filepath.Join("src", "main", "java", filepath.Join(strings.Split(pkg, ".")...), "dto")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			Error("failed to create directories: %v\n", err)
			return
		}

		filePath := filepath.Join(dir, name+"Dto.java")
		content := generatePojoContent(name, pkg, "dto", "Dto", dtoFields)
		if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
			Error("failed to write dto file: %v\n", err)
			return
		}

		Success("Created dto: %s\n", filePath)
	},
}

func init() {
	dtoCmd.Flags().StringVar(&dtoFields, "fields", "", "fields like name:String,age:int")
	dtoCmd.Flags().StringVarP(&dtoPackage, "package", "p", "", "Override base package (ex: com.monentreprise.monprojet)")
	makeCmd.AddCommand(dtoCmd)
}
