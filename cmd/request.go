package cmd

import (
	"fmt"
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
			fmt.Fprintln(os.Stderr, "request name is required")
			return
		}
		name := exportName(raw)

		if !isSpringProject(".") {
			fmt.Fprintln(os.Stderr, "Erreur: Lancez cette commande dans un projet Spring Boot (présence de pom.xml ou build.gradle)")
			os.Exit(1)
		}

		pkg := getEffectivePackage(".", installPackage, requestPackage)

		dir := filepath.Join("src", "main", "java", filepath.Join(strings.Split(pkg, ".")...), "request")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "failed to create directories: %v\n", err)
			return
		}

		filePath := filepath.Join(dir, name+"Request.java")
		sb := strings.Builder{}
		sb.WriteString(fmt.Sprintf("package %s.request;\n\n", pkg))
		sb.WriteString("public class " + name + "Request {\n")

		if strings.TrimSpace(requestFields) != "" {
			parts := strings.Split(requestFields, ",")
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
				sb.WriteString(fmt.Sprintf("    public %s get%s() { return %s; }\n", exportJavaType(ftype), exportName(fname), fname))
				sb.WriteString(fmt.Sprintf("    public void set%s(%s %s) { this.%s = %s; }\n", exportName(fname), exportJavaType(ftype), fname, fname, fname))
			}
		}

		sb.WriteString("}\n")
		if err := os.WriteFile(filePath, []byte(sb.String()), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write request file: %v\n", err)
			return
		}

		fmt.Printf("Created request: %s\n", filePath)
	},
}

func init() {
	requestCmd.Flags().StringVar(&requestFields, "fields", "", "fields like name:String,age:int")
	requestCmd.Flags().StringVarP(&requestPackage, "package", "p", "", "Override base package (ex: com.monentreprise.monprojet)")
	makeCmd.AddCommand(requestCmd)
}
