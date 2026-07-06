package cmd

import (
	"errors"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		raw := args[0]
		if raw == "" {
			return errors.New("service name is required")
		}
		name := exportName(raw)
		if err := requireSpringProject(); err != nil {
			return err
		}

		pkg := getEffectivePackage(".", installPackage, testServicePackage)

		dir := filepath.Join("src", "test", "java", filepath.Join(strings.Split(pkg, ".")...), "service")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("création des dossiers: %w", err)
		}

		file := filepath.Join(dir, name+"ServiceTest.java")
		fields := readEntityFields(pkg, name)
		content, err := renderTemplate("test_service", struct {
			Pkg    string
			Name   string
			Fields []parsedField
		}{Pkg: pkg, Name: name, Fields: fields})
		if err != nil {
			return fmt.Errorf("rendu du template de test: %w", err)
		}
		if err := os.WriteFile(file, []byte(content), 0o644); err != nil {
			return fmt.Errorf("écriture du fichier de test: %w", err)
		}
		Success("Created service test: %s\n", file)
		return nil
	},
}

var testControllerCmd = &cobra.Command{
	Use:   "controller [name]",
	Short: "Generate a JUnit test skeleton for a controller",
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

		pkg := getEffectivePackage(".", installPackage, testControllerPackage)

		dir := filepath.Join("src", "test", "java", filepath.Join(strings.Split(pkg, ".")...), "controller")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("création des dossiers: %w", err)
		}

		file := filepath.Join(dir, name+"ControllerTest.java")
		fields := readEntityFields(pkg, name)
		content, err := renderTemplate("test_controller", struct {
			Pkg    string
			Name   string
			Fields []parsedField
		}{Pkg: pkg, Name: name, Fields: fields})
		if err != nil {
			return fmt.Errorf("rendu du template de test: %w", err)
		}
		if err := os.WriteFile(file, []byte(content), 0o644); err != nil {
			return fmt.Errorf("écriture du fichier de test: %w", err)
		}
		Success("Created controller test: %s\n", file)
		return nil
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
