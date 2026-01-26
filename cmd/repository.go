package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

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

		pkg := "com.example"
		if installPackage != "" {
			pkg = installPackage
		}

		dir := filepath.Join("src", "main", "java", filepath.Join(strings.Split(pkg, ".")...), "repository")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "failed to create directories: %v\n", err)
			return
		}

		filePath := filepath.Join(dir, name+"Repository.java")
		content := fmt.Sprintf("package %s.repository;\n\nimport org.springframework.data.jpa.repository.JpaRepository;\nimport %s.entity.%s;\n\npublic interface %sRepository extends JpaRepository<%s, Long> {\n}\n", pkg, pkg, name, name, name)

		if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write repository file: %v\n", err)
			return
		}

		fmt.Printf("Created repository: %s\n", filePath)
	},
}

func init() {
	makeCmd.AddCommand(repositoryCmd)
}
