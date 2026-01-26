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
	content := fmt.Sprintf("package %s.repository;\n\nimport org.springframework.data.jpa.repository.JpaRepository;\nimport %s.entity.%s;\n\npublic interface %sRepository extends JpaRepository<%s, Long> {\n}\n", pkg, pkg, entity, entity, entity)
	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		return err
	}
	fmt.Printf("Created repository: %s\n", filePath)
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
		iface := fmt.Sprintf("package %s.service;\n\nimport java.util.List;\nimport java.util.Optional;\nimport %s.entity.%s;\n\npublic interface %sService {\n    List<%s> findAll();\n    Optional<%s> findById(Long id);\n    %s save(%s entity);\n    Optional<%s> update(Long id, %s entity);\n    void deleteById(Long id);\n}\n", pkg, pkg, entity, serviceName, entity, entity, entity, entity, entity, entity)
		if err := os.WriteFile(ifacePath, []byte(iface), 0o644); err != nil {
			return err
		}
		fmt.Printf("Created service interface: %s\n", ifacePath)
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

		updateBody := "            // no fields detected to copy\n            return repository.save(existing);\n"
		if copyLines != "" {
			updateBody = copyLines + "\n            return repository.save(existing);\n"
		}

		impl := fmt.Sprintf("package %s.service.impl;\n\nimport java.util.*;\nimport org.springframework.stereotype.Service;\nimport org.springframework.beans.factory.annotation.Autowired;\nimport %s.service.%sService;\nimport %s.entity.%s;\nimport %s.repository.%sRepository;\n\n@Service\npublic class %sServiceImpl implements %sService {\n\n    private final %sRepository repository;\n\n    @Autowired\n    public %sServiceImpl(%sRepository repository) {\n        this.repository = repository;\n    }\n\n    @Override\n    public List<%s> findAll() {\n        return repository.findAll();\n    }\n\n    @Override\n    public Optional<%s> findById(Long id) {\n        return repository.findById(id);\n    }\n\n    @Override\n    public %s save(%s entity) {\n        return repository.save(entity);\n    }\n\n    @Override\n    public Optional<%s> update(Long id, %s entity) {\n        return repository.findById(id).map(existing -> {\n%s    });\n    }\n\n    @Override\n    public void deleteById(Long id) {\n        repository.deleteById(id);\n    }\n}\n", pkg, pkg, serviceName, pkg, entity, pkg, entity, serviceName, serviceName, entity+"Repository", serviceName, entity+"Repository", entity, entity, entity, entity, entity, entity, updateBody)
		if err := os.WriteFile(implPath, []byte(impl), 0o644); err != nil {
			return err
		}
		fmt.Printf("Created service implementation: %s\n", implPath)
	}
	return nil
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
			out = append(out, parsedField{name: n, goType: t})
		}
	}
	return out
}
