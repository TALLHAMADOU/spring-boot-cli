package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Generate test skeletons (JUnit) for services and controllers",
}

var testServiceCmd = &cobra.Command{
	Use:   "service [name]",
	Short: "Generate a JUnit test skeleton for a service",
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

		pkg := getEffectivePackage(".", installPackage, testServicePackage)

		dir := filepath.Join("src", "test", "java", filepath.Join(strings.Split(pkg, ".")...), "service")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "failed to create directories: %v\n", err)
			return
		}

		file := filepath.Join(dir, name+"ServiceTest.java")
		content, err := renderTemplate("test_service", struct {
			Pkg  string
			Name string
		}{Pkg: pkg, Name: name})
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to render test template: %v\n", err)
			return
		}
		if err := os.WriteFile(file, []byte(content), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write test file: %v\n", err)
			return
		}
		fmt.Printf("Created service test: %s\n", file)
	},
}

var testControllerCmd = &cobra.Command{
	Use:   "controller [name]",
	Short: "Generate a JUnit test skeleton for a controller",
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

		pkg := getEffectivePackage(".", installPackage, testControllerPackage)

		dir := filepath.Join("src", "test", "java", filepath.Join(strings.Split(pkg, ".")...), "controller")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "failed to create directories: %v\n", err)
			return
		}

		file := filepath.Join(dir, name+"ControllerTest.java")
		content, err := renderTemplate("test_controller", struct {
			Pkg  string
			Name string
		}{Pkg: pkg, Name: name})
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to render test template: %v\n", err)
			return
		}
		if err := os.WriteFile(file, []byte(content), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write test file: %v\n", err)
			return
		}
		fmt.Printf("Created controller test: %s\n", file)
	},
}

func init() {
	testCmd.AddCommand(testServiceCmd)
	testCmd.AddCommand(testControllerCmd)
	makeCmd.AddCommand(testCmd)
}

var testServicePackage string
var testControllerPackage string

func init() {
	// package override for tests — separate variables to avoid cross-contamination
	testServiceCmd.Flags().StringVarP(&testServicePackage, "package", "p", "", "Override base package for tests (ex: com.example.app)")
	testControllerCmd.Flags().StringVarP(&testControllerPackage, "package", "p", "", "Override base package for tests (ex: com.example.app)")
}
