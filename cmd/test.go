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

		pkg := getEffectivePackage(".", installPackage, testPackage)

		dir := filepath.Join("src", "test", "java", filepath.Join(strings.Split(pkg, ".")...), "service")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "failed to create directories: %v\n", err)
			return
		}

		file := filepath.Join(dir, name+"ServiceTest.java")
		content := fmt.Sprintf("package %s.service;\n\nimport org.junit.jupiter.api.Test;\nimport org.junit.jupiter.api.extension.ExtendWith;\nimport org.mockito.InjectMocks;\nimport org.mockito.Mock;\nimport org.mockito.junit.jupiter.MockitoExtension;\nimport static org.assertj.core.api.Assertions.assertThat;\n\n@ExtendWith(MockitoExtension.class)\npublic class %sServiceTest {\n\n    @Mock\n    // TODO: mock dependencies\n    private Object repository;\n\n    @InjectMocks\n    private %sService service;\n\n    @Test\n    void testFindAll() {\n        // TODO: implement test\n        // assertThat(service.findAll()).isEmpty();\n    }\n}\n", pkg, name, name)
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

		pkg := getEffectivePackage(".", installPackage, testPackage)

		dir := filepath.Join("src", "test", "java", filepath.Join(strings.Split(pkg, ".")...), "controller")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "failed to create directories: %v\n", err)
			return
		}

		file := filepath.Join(dir, name+"ControllerTest.java")
		content := fmt.Sprintf("package %s.controller;\n\nimport org.junit.jupiter.api.Test;\nimport org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;\nimport org.springframework.boot.test.mock.mockito.MockBean;\nimport org.springframework.test.web.servlet.MockMvc;\nimport org.springframework.beans.factory.annotation.Autowired;\nimport static org.assertj.core.api.Assertions.assertThat;\n\n@WebMvcTest(%sController.class)\npublic class %sControllerTest {\n\n    @Autowired\n    private MockMvc mockMvc;\n\n    @MockBean\n    private Object service;\n\n    @Test\n    void testGetAll() throws Exception {\n        // TODO: implement MockMvc test\n        // example assertion: assertThat(mockMvc).isNotNull();\n    }\n}\n", pkg, name, name)
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

var testPackage string

func init() {
	// package override for tests
	testServiceCmd.Flags().StringVarP(&testPackage, "package", "p", "", "Override base package for tests (ex: com.example.app)")
	testControllerCmd.Flags().StringVarP(&testPackage, "package", "p", "", "Override base package for tests (ex: com.example.app)")
}
