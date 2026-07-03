package cmd

import (
	"fmt"
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
		sb := strings.Builder{}
		sb.WriteString(fmt.Sprintf("package %s.dto;\n\n", pkg))
		sb.WriteString("public class " + name + "Dto {\n")

		if strings.TrimSpace(dtoFields) != "" {
			parts := strings.Split(dtoFields, ",")
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p == "" {
					continue
				}
				kv := strings.SplitN(p, ":", 2)
				fname := strings.TrimSpace(kv[0])
				ftype := "String"
				if len(kv) == 2 {
					ftype = strings.TrimSpace(kv[1])
				}
				sb.WriteString(fmt.Sprintf("    private %s %s;\n", exportJavaType(ftype), fname))
				// getter
				sb.WriteString(fmt.Sprintf("    public %s get%s() { return %s; }\n", exportJavaType(ftype), exportName(fname), fname))
				// setter
				sb.WriteString(fmt.Sprintf("    public void set%s(%s %s) { this.%s = %s; }\n", exportName(fname), exportJavaType(ftype), fname, fname, fname))
			}
		}

		sb.WriteString("}\n")
		if err := os.WriteFile(filePath, []byte(sb.String()), 0o644); err != nil {
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
