package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var controllerPackage string

var controllerCmd = &cobra.Command{
	Use:   "controller [name]",
	Short: "Generate a REST controller",
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

		pkg := getEffectivePackage(".", installPackage, controllerPackage)

		dir := filepath.Join("src", "main", "java", filepath.Join(strings.Split(pkg, ".")...), "controller")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "failed to create directories: %v\n", err)
			return
		}

		filePath := filepath.Join(dir, name+"Controller.java")

		crud, _ := cmd.Flags().GetBool("crud")
		entityName, _ := cmd.Flags().GetString("entity")

		var content string
		if crud {
			if entityName == "" {
				fmt.Fprintln(os.Stderr, "--entity is required for --crud")
				return
			}
			e := exportName(entityName)

			// ensure repository and service exist
			if err := ensureRepository(pkg, e); err != nil {
				fmt.Fprintf(os.Stderr, "failed to ensure repository: %v\n", err)
			}
			if err := ensureService(pkg, e+"Service", e); err != nil {
				fmt.Fprintf(os.Stderr, "failed to ensure service: %v\n", err)
			}

			content = fmt.Sprintf("package %s.controller;\n\nimport org.springframework.web.bind.annotation.*;\nimport java.util.List;\nimport java.util.Optional;\nimport %s.entity.%s;\nimport %s.service.%sService;\nimport org.springframework.http.ResponseEntity;\nimport org.springframework.beans.factory.annotation.Autowired;\n\n@RestController\n@RequestMapping(\"/api/%s\")\npublic class %sController {\n\n    @Autowired\n    private %sService service;\n\n    @GetMapping\n    public List<%s> list() {\n        return service.findAll();\n    }\n\n    @GetMapping(\"/{id}\")\n    public ResponseEntity<%s> get(@PathVariable Long id) {\n        return service.findById(id).map(ResponseEntity::ok).orElse(ResponseEntity.notFound().build());\n    }\n\n    @PostMapping\n    public %s create(@RequestBody %s entity) {\n        return service.save(entity);\n    }\n\n    @PutMapping(\"/{id}\")\n    public ResponseEntity<%s> update(@PathVariable Long id, @RequestBody %s entity) {\n        return service.update(id, entity).map(ResponseEntity::ok).orElse(ResponseEntity.notFound().build());\n    }\n\n    @DeleteMapping(\"/{id}\")\n    public ResponseEntity<Void> delete(@PathVariable Long id) {\n        service.deleteById(id);\n        return ResponseEntity.noContent().build();\n    }\n}\n", pkg, pkg, e, pkg, e, strings.ToLower(e), e, e, e, e, e, e, e, e)
		} else {
			content = fmt.Sprintf("package %s.controller;\n\nimport org.springframework.web.bind.annotation.GetMapping;\nimport org.springframework.web.bind.annotation.RequestMapping;\nimport org.springframework.web.bind.annotation.RestController;\n\n@RestController\n@RequestMapping(\"/api/%s\")\npublic class %sController {\n\n    @GetMapping\n    public String index() {\n        return \"ok\";\n    }\n}\n", pkg, strings.ToLower(name), name)
		}

		if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write controller file: %v\n", err)
			return
		}

		fmt.Printf("Created controller: %s\n", filePath)
	},
}

func init() {
	controllerCmd.Flags().Bool("crud", false, "generate CRUD endpoints")
	controllerCmd.Flags().String("entity", "", "entity name to use for CRUD (required with --crud)")
	controllerCmd.Flags().StringVarP(&controllerPackage, "package", "p", "", "Override base package (ex: com.monentreprise.monprojet)")
	makeCmd.AddCommand(controllerCmd)
}

func init() {
	makeCmd.AddCommand(controllerCmd)
}
