package cmd

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// detectVersionInPOM returns the project version (project.version or parent.version)
func detectVersionInPOM(root string) (string, error) {
	pomPath := filepath.Join(root, "pom.xml")
	data, err := os.ReadFile(pomPath)
	if err != nil {
		return "", err
	}
	type parent struct {
		Version string `xml:"version"`
	}
	type project struct {
		Version string  `xml:"version"`
		Parent  *parent `xml:"parent"`
	}
	var p project
	if err := xml.Unmarshal(data, &p); err != nil {
		return "", err
	}
	if strings.TrimSpace(p.Version) != "" {
		return strings.TrimSpace(p.Version), nil
	}
	if p.Parent != nil && strings.TrimSpace(p.Parent.Version) != "" {
		return strings.TrimSpace(p.Parent.Version), nil
	}
	return "", fmt.Errorf("version not found in pom.xml")
}

// setVersionInPOM sets the project-level version (not parent) when possible
func setVersionInPOM(root, newVersion string) error {
	pomPath := filepath.Join(root, "pom.xml")
	data, err := os.ReadFile(pomPath)
	if err != nil {
		return err
	}
	s := string(data)
	// try to find a <version> tag that is not inside <parent>...
	parentIdx := strings.Index(s, "<parent")
	var beforeParent string
	if parentIdx >= 0 {
		beforeParent = s[:parentIdx]
	} else {
		beforeParent = s
	}
	re := regexp.MustCompile(`<version>\s*([^<\s]+)\s*</version>`)
	if m := re.FindStringSubmatch(beforeParent); len(m) >= 2 {
		// replace only the first occurrence in beforeParent
		replaced := re.ReplaceAllString(beforeParent, "<version>"+newVersion+"</version>")
		s = replaced + s[len(beforeParent):]
		return os.WriteFile(pomPath, []byte(s), 0o644)
	}
	// fallback: insert version after <artifactId> if present
	artRe := regexp.MustCompile(`<artifactId>\s*([^<\s]+)\s*</artifactId>`)
	if loc := artRe.FindStringIndex(s); loc != nil {
		insertAt := loc[1]
		insert := "\n  <version>" + newVersion + "</version>\n"
		s = s[:insertAt] + insert + s[insertAt:]
		return os.WriteFile(pomPath, []byte(s), 0o644)
	}
	// last resort: replace first </project> with version section
	s = strings.Replace(s, "</project>", "  <version>"+newVersion+"</version>\n</project>", 1)
	return os.WriteFile(pomPath, []byte(s), 0o644)
}

// detectVersionInGradle finds version = 'x.y.z' or version 'x.y.z'
func detectVersionInGradle(root string) (string, error) {
	gradPath := filepath.Join(root, "build.gradle")
	data, err := os.ReadFile(gradPath)
	if err != nil {
		return "", err
	}
	s := string(data)
	re := regexp.MustCompile(`(?m)^[ \t]*version\s*(?:=\s*)?['\"]([^'\"]+)['\"]`)
	if m := re.FindStringSubmatch(s); len(m) >= 2 {
		return strings.TrimSpace(m[1]), nil
	}
	return "", fmt.Errorf("version not found in build.gradle")
}

// setVersionInGradle replaces the version line keeping rest intact
func setVersionInGradle(root, newVersion string) error {
	gradPath := filepath.Join(root, "build.gradle")
	data, err := os.ReadFile(gradPath)
	if err != nil {
		return err
	}
	s := string(data)
	re := regexp.MustCompile(`(?m)^(?P<prefix>[ \t]*version\s*(?:=\s*)?)[ '\"]?[^'\"\n]+[ '\"]?`)
	if re.MatchString(s) {
		s = re.ReplaceAllString(s, "${prefix}'"+newVersion+"'\n")
		return os.WriteFile(gradPath, []byte(s), 0o644)
	}
	// fallback: append version = 'x'
	s = s + "\nversion = '" + newVersion + "'\n"
	return os.WriteFile(gradPath, []byte(s), 0o644)
}

// detectVersionInBuildFiles returns version and kind (maven|gradle)
func detectVersionInBuildFiles(root string) (string, string, error) {
	if v, err := detectVersionInPOM(root); err == nil {
		return v, "maven", nil
	}
	if v, err := detectVersionInGradle(root); err == nil {
		return v, "gradle", nil
	}
	return "", "", fmt.Errorf("no version detected in pom.xml or build.gradle")
}

// bumpSemver increments a semver string (simple x.y.z)
func bumpSemver(v, kind string) (string, error) {
	orig := strings.TrimSpace(v)
	orig = strings.TrimPrefix(orig, "v")
	parts := strings.Split(orig, ".")
	for len(parts) < 3 {
		parts = append(parts, "0")
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return "", err
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", err
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", err
	}
	switch kind {
	case "patch":
		patch++
	case "minor":
		minor++
		patch = 0
	case "major":
		major++
		minor = 0
		patch = 0
	default:
		return "", fmt.Errorf("unknown bump kind: %s", kind)
	}
	return fmt.Sprintf("%d.%d.%d", major, minor, patch), nil
}
