package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/spf13/cobra"
)

var (
	entityFields   string
	entityLombok   bool
	entityAuditing  bool
	entityPackage   string
	entityUUID      bool
	entityValidate  bool
	entityHasMany   string
	entityBelongsTo string
	entityHasOne    string
)

var entityCmd = &cobra.Command{
	Use:   "entity [name]",
	Short: "Generate a JPA entity class",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		raw := args[0]
		if raw == "" {
			return errors.New("entity name is required")
		}
		name := exportName(raw)

		// ensure we are in a Spring Boot project
		if err := requireSpringProject(); err != nil {
			return err
		}

		// detect base package from project files (pom.xml / build.gradle)
		detected := detectBasePackage(".")
		pkg := detected
		// installPackage (from install command) has priority over detected
		if installPackage != "" {
			pkg = installPackage
		}
		// explicit --package flag for this command overrides everything
		if entityPackage != "" {
			pkg = entityPackage
		}

		dir := filepath.Join("src", "main", "java", filepath.Join(strings.Split(pkg, ".")...), "entity")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("création des dossiers: %w", err)
		}

		filePath := filepath.Join(dir, name+".java")

		imports := []string{"jakarta.persistence.Entity", "jakarta.persistence.Id", "jakarta.persistence.GeneratedValue", "jakarta.persistence.GenerationType"}

		if entityAuditing {
			imports = append(imports,
				"jakarta.persistence.EntityListeners",
				"org.springframework.data.annotation.CreatedDate",
				"org.springframework.data.annotation.LastModifiedDate",
				"org.springframework.data.jpa.domain.support.AuditingEntityListener",
			)
		}

		if entityValidate {
			imports = append(imports, "jakarta.validation.constraints.NotBlank", "jakarta.validation.constraints.NotNull")
		}

		if entityLombok {
			imports = append(imports, "lombok.Getter", "lombok.Setter", "lombok.NoArgsConstructor", "lombok.AllArgsConstructor")
		}

		// parse additional fields
		var fields []entityField
		if strings.TrimSpace(entityFields) != "" {
			for _, f := range parseFields(entityFields) {
				imports = append(imports, f.importPkg...)
				fields = append(fields, entityField{Type: f.Type, Name: f.Name, Cap: exportName(f.Name)})
			}
		}

		hasManyFields := parseRelation(entityHasMany, "List<%s>", "s")
		belongsToFields := parseRelation(entityBelongsTo, "%s", "")
		hasOneFields := parseRelation(entityHasOne, "%s", "")

		if len(hasManyFields) > 0 {
			imports = append(imports, "jakarta.persistence.OneToMany", "java.util.List")
		}
		if len(belongsToFields) > 0 {
			imports = append(imports, "jakarta.persistence.ManyToOne")
		}
		if len(hasOneFields) > 0 {
			imports = append(imports, "jakarta.persistence.OneToOne")
		}

		// remove duplicate imports
		imports = uniqueStrings(imports)

		content, err := renderTemplate("entity", entityData{
			Pkg:       pkg,
			Name:      name,
			Imports:   imports,
			Auditing:  entityAuditing,
			Lombok:    entityLombok,
			UUID:      entityUUID,
			Validate:  entityValidate,
			Fields:    fields,
			HasMany:   hasManyFields,
			BelongsTo: belongsToFields,
			HasOne:    hasOneFields,
		})
		if err != nil {
			return fmt.Errorf("rendu du template entity: %w", err)
		}

		if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
			return fmt.Errorf("écriture du fichier entity: %w", err)
		}

		// If lombok requested, try to add dependency to project's pom.xml or build.gradle in current working dir
		if entityLombok {
			addLombokToProject(".")
		}

		// ensure JPA dependency exists for entity generation (if pom/build.gradle present)
		if err := ensureJPAInProject("."); err != nil {
			Warning("failed to ensure JPA dependency: %v\n", err)
		}

		Success("Created entity: %s\n", filePath)
		return nil
	},
}

func init() {
	entityCmd.Flags().StringVar(&entityFields, "fields", "", "comma-separated fields like name:String,age:int")
	entityCmd.Flags().BoolVar(&entityLombok, "lombok", false, "use Lombok annotations instead of getters/setters")
	entityCmd.Flags().BoolVar(&entityAuditing, "auditing", false, "add createdAt/updatedAt auditing fields")
	entityCmd.Flags().StringVar(&entityPackage, "package", "", "base package override (e.g. com.example.app)")
	entityCmd.Flags().BoolVar(&entityUUID, "uuid", false, "use UUID instead of Long for the primary key")
	entityCmd.Flags().BoolVar(&entityValidate, "validate", false, "add Jakarta validation annotations (@NotBlank, @NotNull)")
	entityCmd.Flags().StringVar(&entityHasMany, "has-many", "", "comma-separated entities (e.g. Order,Item)")
	entityCmd.Flags().StringVar(&entityBelongsTo, "belongs-to", "", "comma-separated entities (e.g. User,Category)")
	entityCmd.Flags().StringVar(&entityHasOne, "has-one", "", "comma-separated entities (e.g. Profile)")
	makeCmd.AddCommand(entityCmd)
}

func exportName(s string) string {
	r := []rune(s)
	if len(r) == 0 {
		return s
	}
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

type parsedField struct {
	Name      string
	Type      string
	importPkg []string
}

// entityField is a single field rendered by the entity template.
type entityField struct {
	Type string
	Name string
	Cap  string // exported name used for getter/setter identifiers
}

// entityData is the data model passed to the entity template.
type entityData struct {
	Pkg      string
	Name     string
	Imports  []string
	Auditing bool
	Lombok   bool
	UUID      bool
	Validate  bool
	Fields    []entityField
	HasMany   []relationField
	BelongsTo []relationField
	HasOne    []relationField
}

type relationField struct {
	Type   string
	Entity string
	Name   string
	Cap    string
}

func parseRelation(spec, typeFormat, nameSuffix string) []relationField {
	parts := strings.Split(spec, ",")
	out := make([]relationField, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		entityName := exportName(p)
		fieldName := strings.ToLower(string(p[0])) + p[1:] + nameSuffix
		fieldType := fmt.Sprintf(typeFormat, entityName)
		out = append(out, relationField{Type: fieldType, Entity: entityName, Name: fieldName, Cap: exportName(fieldName)})
	}
	return out
}

func parseFields(spec string) []parsedField {
	parts := strings.Split(spec, ",")
	out := make([]parsedField, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		kv := strings.SplitN(p, ":", 2)
		fname := strings.TrimSpace(kv[0])
		ftype := "String"
		if len(kv) == 2 {
			switch strings.ToLower(strings.TrimSpace(kv[1])) {
			case "string", "str":
				ftype = "String"
			case "int", "integer":
				ftype = "Integer"
			case "long":
				ftype = "Long"
			case "boolean", "bool":
				ftype = "Boolean"
			case "double":
				ftype = "Double"
			case "instant", "datetime":
				ftype = "java.time.Instant"
			default:
				ftype = exportJavaType(kv[1])
			}
		}
		// java.time fields are emitted with fully-qualified names in the body
		// (like the auditing fields), so no separate import entry is needed.
		out = append(out, parsedField{Name: fname, Type: ftype})
	}
	return out
}

func exportJavaType(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "String"
	}
	// naive: if contains dot, assume full type
	if strings.Contains(s, ".") {
		return s
	}
	// capitalize
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func uniqueStrings(in []string) []string {
	seen := map[string]struct{}{}
	out := []string{}
	for _, s := range in {
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

// addLombokToProject tries to insert Lombok dependency into pom.xml or build.gradle
func addLombokToProject(root string) {
	pomPath := filepath.Join(root, "pom.xml")
	if _, err := os.Stat(pomPath); err == nil {
		dep := `    <dependency>
			<groupId>org.projectlombok</groupId>
			<artifactId>lombok</artifactId>
			<version>1.18.26</version>
			<scope>provided</scope>
		</dependency>
`
		if err := insertPOMDependency(pomPath, "org.projectlombok", dep); err != nil {
			Warning("failed to add Lombok dependency: %v\n", err)
			return
		}
		plugin := `      <plugin>
				<groupId>org.apache.maven.plugins</groupId>
				<artifactId>maven-compiler-plugin</artifactId>
				<version>3.10.1</version>
				<configuration>
					<annotationProcessorPaths>
						<path>
							<groupId>org.projectlombok</groupId>
							<artifactId>lombok</artifactId>
							<version>1.18.26</version>
						</path>
					</annotationProcessorPaths>
				</configuration>
			</plugin>
`
		if err := insertPOMPlugin(pomPath, "maven-compiler-plugin", plugin); err != nil {
			Warning("failed to add Lombok compiler plugin: %v\n", err)
		}
		Success("Added Lombok dependency and compiler plugin to %s\n", pomPath)
		return
	}

	gradlePath := filepath.Join(root, "build.gradle")
	if _, err := os.Stat(gradlePath); err == nil {
		dep := "    compileOnly 'org.projectlombok:lombok:1.18.26'\n    annotationProcessor 'org.projectlombok:lombok:1.18.26'\n"
		if err := insertGradleDependency(gradlePath, "org.projectlombok", dep); err != nil {
			Warning("failed to add Lombok dependency: %v\n", err)
			return
		}
		Success("Added Lombok dependency to %s\n", gradlePath)
	}
}
