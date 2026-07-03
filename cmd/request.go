package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var requestFields string
var requestPackage string

var requestCmd = &cobra.Command{
	Use:   "request [name]",
	Short: "Generate a Request class (for controller inputs)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		raw := args[0]
		if raw == "" {
			Error("request name is required")
			return
		}
		name := exportName(raw)

		if !isSpringProject(".") {
			Error("Erreur: Lancez cette commande dans un projet Spring Boot (présence de pom.xml ou build.gradle)")
			os.Exit(1)
		}

		pkg := getEffectivePackage(".", installPackage, requestPackage)

		dir := filepath.Join("src", "main", "java", filepath.Join(strings.Split(pkg, ".")...), "request")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			Error("failed to create directories: %v\n", err)
			return
		}

		filePath := filepath.Join(dir, name+"Request.java")
		content := generatePojoContent(name, pkg, "request", "Request", requestFields)
		if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
			Error("failed to write request file: %v\n", err)
			return
		}

		Success("Created request: %s\n", filePath)
	},
}

func init() {
	requestCmd.Flags().StringVar(&requestFields, "fields", "", "fields like name:String,age:int")
	requestCmd.Flags().StringVarP(&requestPackage, "package", "p", "", "Override base package (ex: com.monentreprise.monprojet)")
	makeCmd.AddCommand(requestCmd)
}
