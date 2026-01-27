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
			if gid != "" {
				base := strings.ToLower(gid)
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
			aid = strings.TrimSpace(all[len(all)-1][1])
		}
		if gid != "" {
			return strings.ToLower(gid)
		}
		if aid != "" {
			return strings.ToLower(strings.ReplaceAll(aid, "-", ""))
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

// ensureJPAInProject ensures that spring-boot-starter-data-jpa is present in pom.xml if it exists
func ensureJPAInProject(root string) error {
	pomPath := filepath.Join(root, "pom.xml")
	data, err := os.ReadFile(pomPath)
	if err != nil {
		// no pom.xml, try build.gradle
		gradPath := filepath.Join(root, "build.gradle")
		g, err2 := os.ReadFile(gradPath)
		if err2 != nil {
			return nil // nothing to do
		}
		s := string(g)
		if strings.Contains(s, "spring-boot-starter-data-jpa") {
			return nil
		}
		add := "\n// Data JPA\nimplementation 'org.springframework.boot:spring-boot-starter-data-jpa'\n"
		if strings.Contains(s, "dependencies {") {
			s = strings.Replace(s, "dependencies {", "dependencies {"+add, 1)
		} else {
			s = s + "\ndependencies {" + add + "}\n"
		}
		return os.WriteFile(gradPath, []byte(s), 0o644)
	}

	s := string(data)
	if strings.Contains(s, "spring-boot-starter-data-jpa") {
		return nil
	}

	dep := `    <dependency>
      <groupId>org.springframework.boot</groupId>
      <artifactId>spring-boot-starter-data-jpa</artifactId>
    </dependency>
`
	if strings.Contains(s, "</dependencies>") {
		s = strings.Replace(s, "</dependencies>", dep+"  </dependencies>", 1)
	} else {
		s = strings.Replace(s, "</project>", "  <dependencies>\n"+dep+"  </dependencies>\n</project>", 1)
	}
	if err := os.WriteFile(pomPath, []byte(s), 0o644); err != nil {
		return err
	}
	fmt.Printf("Added Spring Data JPA dependency to %s\n", pomPath)
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
