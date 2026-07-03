package cmd

import (
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
	entityAuditing bool
	entityPackage  string
)

var entityCmd = &cobra.Command{
	Use:   "entity [name]",
	Short: "Generate a JPA entity class",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		raw := args[0]
		if raw == "" {
			Error("entity name is required")
			return
		}
		name := exportName(raw)

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

		// ensure we are in a Spring Boot project
		if !isSpringProject(".") {
			Error("Erreur: Lancez dans un projet Spring Boot avec pom.xml ou build.gradle")
			os.Exit(1)
		}

		dir := filepath.Join("src", "main", "java", filepath.Join(strings.Split(pkg, ".")...), "entity")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			Error("failed to create directories: %v\n", err)
			return
		}

		filePath := filepath.Join(dir, name+".java")

		imports := []string{"jakarta.persistence.Entity", "jakarta.persistence.Id", "jakarta.persistence.GeneratedValue", "jakarta.persistence.GenerationType"}
		bodyFields := []string{"    @Id", "    @GeneratedValue(strategy = GenerationType.IDENTITY)", "    private Long id;"}
		getters := []string{"    public Long getId() {", "        return id;", "    }"}
		setters := []string{"    public void setId(Long id) {", "        this.id = id;", "    }"}

		if entityAuditing {
			imports = append(imports,
				"jakarta.persistence.EntityListeners",
				"org.springframework.data.annotation.CreatedDate",
				"org.springframework.data.annotation.LastModifiedDate",
				"org.springframework.data.jpa.domain.support.AuditingEntityListener",
			)
			bodyFields = append(bodyFields, "", "    // Auditing fields", "    @CreatedDate", "    private java.time.Instant createdAt;", "", "    @LastModifiedDate", "    private java.time.Instant updatedAt;")
		}

		if entityLombok {
			imports = append(imports, "lombok.Getter", "lombok.Setter", "lombok.NoArgsConstructor", "lombok.AllArgsConstructor")
		}

		// parse additional fields
		if strings.TrimSpace(entityFields) != "" {
			parsed := parseFields(entityFields)
			for _, f := range parsed {
				imports = append(imports, f.importPkg...)
				bodyFields = append(bodyFields, "", fmt.Sprintf("    private %s %s;", f.goType, f.name))
				if !entityLombok {
					// add getter
					getters = append(getters, "", fmt.Sprintf("    public %s get%s() {", f.goType, exportName(f.name)), fmt.Sprintf("        return %s;", f.name), "    }")
					// add setter
					setters = append(setters, "", fmt.Sprintf("    public void set%s(%s %s) {", exportName(f.name), f.goType, f.name), fmt.Sprintf("        this.%s = %s;", f.name, f.name), "    }")
				}
			}
		}

		// remove duplicate imports
		imports = uniqueStrings(imports)

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("package %s.entity;\n\n", pkg))
		for _, im := range imports {
			sb.WriteString(fmt.Sprintf("import %s;\n", im))
		}
		sb.WriteString("\n")
		if entityAuditing {
			sb.WriteString("@EntityListeners(AuditingEntityListener.class)\n")
		}
		sb.WriteString("@Entity\n")
		sb.WriteString(fmt.Sprintf("public class %s {\n", name))
		for _, l := range bodyFields {
			sb.WriteString(l + "\n")
		}
		sb.WriteString("\n")
		if !entityLombok {
			for _, g := range getters {
				sb.WriteString(g + "\n")
			}
			sb.WriteString("\n")
			for _, s := range setters {
				sb.WriteString(s + "\n")
			}
		} else {
			// Lombok annotations
			sb.WriteString("\n")
			sb.WriteString("@Getter\n@Setter\n@NoArgsConstructor\n@AllArgsConstructor\n")
		}
		sb.WriteString("}\n")

		if err := os.WriteFile(filePath, []byte(sb.String()), 0o644); err != nil {
			Error("failed to write entity file: %v\n", err)
			return
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
	},
}

func init() {
	entityCmd.Flags().StringVar(&entityFields, "fields", "", "comma-separated fields like name:String,age:int")
	entityCmd.Flags().BoolVar(&entityLombok, "lombok", false, "use Lombok annotations instead of getters/setters")
	entityCmd.Flags().BoolVar(&entityAuditing, "auditing", false, "add createdAt/updatedAt auditing fields")
	entityCmd.Flags().StringVar(&entityPackage, "package", "", "base package override (e.g. com.example.app)")
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
	name      string
	goType    string
	importPkg []string
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
		imp := []string{}
		if strings.Contains(ftype, "java.time") {
			imp = append(imp, "java.time.Instant")
		}
		out = append(out, parsedField{name: fname, goType: ftype, importPkg: imp})
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
