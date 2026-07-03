package cmd

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// isSpringProject returns true if pom.xml or build.gradle exists in root
func isSpringProject(root string) bool {
	if _, err := os.Stat(filepath.Join(root, "pom.xml")); err == nil {
		return true
	}
	if _, err := os.Stat(filepath.Join(root, "build.gradle")); err == nil {
		return true
	}
	return false
}

// detectBasePackage tries to detect the base java package for the project located at root
// Priority: pom.xml (groupId + artifactId) -> build.gradle + settings.gradle -> fallback com.example
func detectBasePackage(root string) string {
	// try pom.xml
	pomPath := filepath.Join(root, "pom.xml")
	if data, err := os.ReadFile(pomPath); err == nil {
		type parent struct {
			GroupId    string `xml:"groupId"`
			ArtifactId string `xml:"artifactId"`
		}
		type project struct {
			XMLName    xml.Name `xml:"project"`
			GroupId    string   `xml:"groupId"`
			ArtifactId string   `xml:"artifactId"`
			Parent     *parent  `xml:"parent"`
		}
		var p project
		if err := xml.Unmarshal(data, &p); err == nil {
			gid := strings.TrimSpace(p.GroupId)
			if gid == "" && p.Parent != nil {
				gid = strings.TrimSpace(p.Parent.GroupId)
			}
			aid := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(p.ArtifactId), "-", ""))
			if gid != "" {
				base := strings.ToLower(gid)
				if aid != "" {
					return base + "." + aid
				}
				return base
			}
		}

		// fallback: try tolerant regex-based extraction (handles default namespace and parent blocks)
		s := string(data)
		// limit search to the header section (before <dependencies>) to avoid plugin/groupId occurrences
		depIdx := strings.Index(s, "<dependencies")
		header := s
		if depIdx >= 0 {
			header = s[:depIdx]
		}
		// remove parent block from header if present
		if strings.Contains(header, "<parent") {
			start := strings.Index(header, "<parent")
			end := strings.Index(header, "</parent>")
			if start >= 0 && end > start {
				header = header[:start] + header[end+len("</parent>"):]
			}
		}
		reG := regexp.MustCompile(`<groupId>\s*([^<\s]+)\s*</groupId>`)
		reA := regexp.MustCompile(`<artifactId>\s*([^<\s]+)\s*</artifactId>`)
		gid := ""
		aid := ""
		// find last groupId/artifactId in header
		if all := reG.FindAllStringSubmatch(header, -1); len(all) > 0 {
			gid = strings.TrimSpace(all[len(all)-1][1])
		}
		if all := reA.FindAllStringSubmatch(header, -1); len(all) > 0 {
			aid = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(all[len(all)-1][1]), "-", ""))
		}
		if gid != "" && aid != "" {
			return strings.ToLower(gid) + "." + aid
		}
		if gid != "" {
			return strings.ToLower(gid)
		}
		if aid != "" {
			return strings.ToLower(aid)
		}
	}

	// try build.gradle + settings.gradle
	gradPath := filepath.Join(root, "build.gradle")
	if data, err := os.ReadFile(gradPath); err == nil {
		s := string(data)
		// look for group = 'com.example'
		re := regexp.MustCompile(`group\s*=\s*['\"]([^'\"]+)['\"]`)
		grp := ""
		if m := re.FindStringSubmatch(s); len(m) >= 2 {
			grp = strings.TrimSpace(m[1])
		}
		// settings.gradle for project name (not used for base pkg unless group unknown)
		if grp != "" {
			return strings.ToLower(grp)
		}
	}

	// fallback
	fmt.Fprintf(os.Stderr, "Warning: failed to detect base package, falling back to com.example\n")
	return "com.example"
}

// insertPOMDependency safely inserts a dependency XML block into pom.xml.
// sentinel: unique substring identifying the dependency (used to avoid duplicates).
// depXML: the full <dependency>...</dependency> block including leading whitespace.
func insertPOMDependency(pomPath, sentinel, depXML string) error {
	data, err := os.ReadFile(pomPath)
	if err != nil {
		return err
	}
	s := string(data)
	if strings.Contains(s, sentinel) {
		return nil // already present
	}
	// Locate <dependencies> opening tag, then find the matching </dependencies>.
	// This avoids accidentally matching </dependencies> inside a <parent> block.
	depOpenRe := regexp.MustCompile(`<dependencies\b[^>]*>`)
	if loc := depOpenRe.FindStringIndex(s); loc != nil {
		rest := s[loc[1]:]
		if closeIdx := strings.Index(rest, "</dependencies>"); closeIdx >= 0 {
			insertAt := loc[1] + closeIdx
			s = s[:insertAt] + depXML + s[insertAt:]
			return os.WriteFile(pomPath, []byte(s), 0o644)
		}
	}
	// No <dependencies> section: create one before </project>.
	if idx := strings.LastIndex(s, "</project>"); idx >= 0 {
		s = s[:idx] + "  <dependencies>\n" + depXML + "  </dependencies>\n" + s[idx:]
		return os.WriteFile(pomPath, []byte(s), 0o644)
	}
	return fmt.Errorf("cannot insert dependency into %s: no suitable insertion point", pomPath)
}

// insertPOMPlugin safely inserts a plugin XML block inside <build><plugins>.
func insertPOMPlugin(pomPath, sentinel, pluginXML string) error {
	data, err := os.ReadFile(pomPath)
	if err != nil {
		return err
	}
	s := string(data)
	if strings.Contains(s, sentinel) {
		return nil
	}
	pluginOpenRe := regexp.MustCompile(`<plugins\b[^>]*>`)
	if loc := pluginOpenRe.FindStringIndex(s); loc != nil {
		rest := s[loc[1]:]
		if closeIdx := strings.Index(rest, "</plugins>"); closeIdx >= 0 {
			insertAt := loc[1] + closeIdx
			s = s[:insertAt] + pluginXML + s[insertAt:]
			return os.WriteFile(pomPath, []byte(s), 0o644)
		}
	}
	// Fallback: add inside <build>, or before </project>.
	if strings.Contains(s, "</build>") {
		s = strings.Replace(s, "</build>", "  <plugins>\n"+pluginXML+"  </plugins>\n</build>", 1)
	} else if idx := strings.LastIndex(s, "</project>"); idx >= 0 {
		s = s[:idx] + "  <build>\n    <plugins>\n" + pluginXML + "    </plugins>\n  </build>\n" + s[idx:]
	}
	return os.WriteFile(pomPath, []byte(s), 0o644)
}

// insertGradleDependency safely inserts a dependency line into build.gradle.
func insertGradleDependency(gradlePath, sentinel, depLine string) error {
	data, err := os.ReadFile(gradlePath)
	if err != nil {
		return err
	}
	s := string(data)
	if strings.Contains(s, sentinel) {
		return nil
	}
	if strings.Contains(s, "dependencies {") {
		s = strings.Replace(s, "dependencies {", "dependencies {\n"+depLine, 1)
	} else {
		s = s + "\ndependencies {\n" + depLine + "}\n"
	}
	return os.WriteFile(gradlePath, []byte(s), 0o644)
}

// ensureJPAInProject ensures that spring-boot-starter-data-jpa is present in pom.xml or build.gradle.
func ensureJPAInProject(root string) error {
	pomPath := filepath.Join(root, "pom.xml")
	if _, err := os.Stat(pomPath); err == nil {
		dep := `    <dependency>
      <groupId>org.springframework.boot</groupId>
      <artifactId>spring-boot-starter-data-jpa</artifactId>
    </dependency>
`
		if err := insertPOMDependency(pomPath, "spring-boot-starter-data-jpa", dep); err != nil {
			return err
		}
		fmt.Printf("Added Spring Data JPA dependency to %s\n", pomPath)
		return nil
	}
	gradlePath := filepath.Join(root, "build.gradle")
	if _, err := os.Stat(gradlePath); err == nil {
		dep := "    implementation 'org.springframework.boot:spring-boot-starter-data-jpa'\n"
		if err := insertGradleDependency(gradlePath, "spring-boot-starter-data-jpa", dep); err != nil {
			return err
		}
		fmt.Printf("Added Spring Data JPA dependency to %s\n", gradlePath)
	}
	return nil
}

// getEffectivePackage applies detection and overrides: detected <- installPkg <- override
func getEffectivePackage(root, installPkg, override string) string {
	base := detectBasePackage(root)
	pkg := base
	if strings.TrimSpace(installPkg) != "" {
		pkg = installPkg
	}
	if strings.TrimSpace(override) != "" {
		pkg = override
	}
	return pkg
}
