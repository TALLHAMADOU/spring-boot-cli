package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ensureRepository creates a repository for the given entity if it doesn't exist
func ensureRepository(pkg, entity string) error {
	dir := filepath.Join("src", "main", "java", filepath.Join(strings.Split(pkg, ".")...), "repository")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	filePath := filepath.Join(dir, entity+"Repository.java")
	if _, err := os.Stat(filePath); err == nil {
		return nil // exists
	}
	content, err := renderTemplate("repository", struct {
		Pkg    string
		Entity string
	}{Pkg: pkg, Entity: entity})
	if err != nil {
		return fmt.Errorf("render repository template: %w", err)
	}
	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		return err
	}
	Success("Created repository: %s\n", filePath)
	return nil
}

// ensureService creates a service interface and implementation for the given entity if missing
func ensureService(pkg, serviceName, entity string) error {
	// ensure repository exists for the entity so the service can use it
	if entity != "" {
		if err := ensureRepository(pkg, entity); err != nil {
			return err
		}
	}
	dir := filepath.Join("src", "main", "java", filepath.Join(strings.Split(pkg, ".")...), "service")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	ifacePath := filepath.Join(dir, serviceName+"Service.java")
	implDir := filepath.Join(dir, "impl")
	if err := os.MkdirAll(implDir, 0o755); err != nil {
		return err
	}
	implPath := filepath.Join(implDir, serviceName+"ServiceImpl.java")

	if _, err := os.Stat(ifacePath); os.IsNotExist(err) {
		iface, err := renderTemplate("service_interface", struct {
			Pkg    string
			Name   string
			Entity string
		}{Pkg: pkg, Name: serviceName, Entity: entity})
		if err != nil {
			return fmt.Errorf("render service_interface template: %w", err)
		}
		if err := os.WriteFile(ifacePath, []byte(iface), 0o644); err != nil {
			return err
		}
		Success("Created service interface: %s\n", ifacePath)
	}

	if _, err := os.Stat(implPath); os.IsNotExist(err) {
		// attempt to read entity fields to generate explicit copy statements
		fields := readEntityFields(pkg, entity)
		copyLines := ""
		for _, f := range fields {
			if f.name == "id" || f.name == "createdAt" || f.name == "updatedAt" {
				continue
			}
			cap := exportName(f.name)
			copyLines += fmt.Sprintf("            existing.set%s(entity.get%s());\n", cap, cap)
		}
		if copyLines == "" {
			copyLines = "            // no fields detected to copy\n            return repository.save(existing);\n"
		} else {
			copyLines += "\n            return repository.save(existing);\n"
		}
		impl, err := renderTemplate("service_impl", struct {
			Pkg       string
			Name      string
			Entity    string
			CopyLines string
		}{Pkg: pkg, Name: serviceName, Entity: entity, CopyLines: copyLines})
		if err != nil {
			return fmt.Errorf("render service_impl template: %w", err)
		}
		if err := os.WriteFile(implPath, []byte(impl), 0o644); err != nil {
			return err
		}
		Success("Created service implementation: %s\n", implPath)
	}
	return nil
}

// generatePojoContent generates a Java POJO class (DTO, Request, etc.) content string.
// name: capitalized class name, pkg: base package, subPkg: sub-package (e.g. "dto"),
// suffix: class name suffix (e.g. "Dto", "Request"), fieldsSpec: "name:Type,..." pairs.
func generatePojoContent(name, pkg, subPkg, suffix, fieldsSpec string, validate bool) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "package %s.%s;\n\n", pkg, subPkg)
	if validate {
		sb.WriteString("import jakarta.validation.constraints.NotBlank;\n")
		sb.WriteString("import jakarta.validation.constraints.NotNull;\n\n")
	}
	fmt.Fprintf(&sb, "public class %s%s {\n", name, suffix)
	if strings.TrimSpace(fieldsSpec) != "" {
		for p := range strings.SplitSeq(fieldsSpec, ",") {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			kv := strings.SplitN(p, ":", 2)
			fname := strings.TrimSpace(kv[0])
			ftype := "String"
			if len(kv) == 2 {
				ftype = exportJavaType(strings.TrimSpace(kv[1]))
			}
			if validate {
				if ftype == "String" {
					fmt.Fprintf(&sb, "    @NotBlank\n")
				} else {
					fmt.Fprintf(&sb, "    @NotNull\n")
				}
			}
			fmt.Fprintf(&sb, "    private %s %s;\n", ftype, fname)
			fmt.Fprintf(&sb, "    public %s get%s() { return %s; }\n", ftype, exportName(fname), fname)
			fmt.Fprintf(&sb, "    public void set%s(%s %s) { this.%s = %s; }\n",
				exportName(fname), ftype, fname, fname, fname)
		}
	}
	sb.WriteString("}\n")
	return sb.String()
}

// readEntityFields attempts to parse private fields from the generated entity java file
func readEntityFields(pkg, entity string) []parsedField {
	path := filepath.Join("src", "main", "java", filepath.Join(strings.Split(pkg, ".")...), "entity", entity+".java")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	re := regexp.MustCompile(`private\s+([A-Za-z0-9_\.<>]+)\s+([A-Za-z0-9_]+);`)
	matches := re.FindAllStringSubmatch(string(data), -1)
	out := []parsedField{}
	for _, m := range matches {
		if len(m) >= 3 {
			t := m[1]
			n := m[2]
			// ignore synthetic/auditing or id types handled elsewhere
			out = append(out, parsedField{Name: n, Type: t})
		}
	}
	return out
}
