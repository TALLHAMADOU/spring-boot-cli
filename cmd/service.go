package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var servicePackage string

var serviceCmd = &cobra.Command{
	Use:   "service [name]",
	Short: "Generate a service interface and implementation",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		raw := args[0]
		if raw == "" {
			fmt.Fprintln(os.Stderr, "service name is required")
			return
		}
		name := exportName(raw)
		if !isSpringProject(".") {
			fmt.Fprintln(os.Stderr, "Erreur: Lancez cette commande dans un projet Spring Boot (présence de pom.xml ou build.gradle)")
			os.Exit(1)
		}

		pkg := getEffectivePackage(".", installPackage, servicePackage)

		entityName, _ := cmd.Flags().GetString("entity")

		if entityName != "" {
			e := exportName(entityName)
			if err := ensureService(pkg, name, e); err != nil {
				fmt.Fprintf(os.Stderr, "failed to ensure service: %v\n", err)
				return
			}
			fmt.Printf("Ensured service and implementation for %s (entity %s)\n", name, e)
			return
		}

		// non-entity generic service
		dir := filepath.Join("src", "main", "java", filepath.Join(strings.Split(pkg, ".")...), "service")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "failed to create directories: %v\n", err)
			return
		}

		ifacePath := filepath.Join(dir, name+"Service.java")
		iface, err := renderTemplate("service_interface_generic", struct {
			Pkg  string
			Name string
		}{Pkg: pkg, Name: name})
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to render service interface template: %v\n", err)
			return
		}
		if err := os.WriteFile(ifacePath, []byte(iface), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write service interface: %v\n", err)
			return
		}

		implPath := filepath.Join(dir, "impl")
		if err := os.MkdirAll(implPath, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "failed to create impl dir: %v\n", err)
			return
		}
		implFile := filepath.Join(implPath, name+"ServiceImpl.java")
		impl, err := renderTemplate("service_impl_generic", struct {
			Pkg  string
			Name string
		}{Pkg: pkg, Name: name})
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to render service impl template: %v\n", err)
			return
		}
		if err := os.WriteFile(implFile, []byte(impl), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write service implementation: %v\n", err)
			return
		}

		fmt.Printf("Created service: %s and implementation %s\n", ifacePath, implFile)
	},
}

func init() {
	makeCmd.AddCommand(serviceCmd)
	serviceCmd.Flags().StringP("entity", "e", "", "associate service with an entity (generate entity-backed methods)")
	serviceCmd.Flags().StringVarP(&servicePackage, "package", "p", "", "Override base package (ex: com.monentreprise.monprojet)")
}
