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

// parentBlockRange returns the [start, end) byte range of the <parent>...</parent>
// block in s, or (-1, -1) if there is none.
func parentBlockRange(s string) (int, int) {
	start := strings.Index(s, "<parent")
	if start < 0 {
		return -1, -1
	}
	closeIdx := strings.Index(s[start:], "</parent>")
	if closeIdx < 0 {
		return -1, -1
	}
	return start, start + closeIdx + len("</parent>")
}

// outsideParent reports whether byte offset idx falls outside the parent block.
func outsideParent(idx, pStart, pEnd int) bool {
	return pStart < 0 || idx < pStart || idx >= pEnd
}

// setVersionInPOM sets the project-level version (not the parent version).
func setVersionInPOM(root, newVersion string) error {
	pomPath := filepath.Join(root, "pom.xml")
	data, err := os.ReadFile(pomPath)
	if err != nil {
		return err
	}
	s := string(data)
	pStart, pEnd := parentBlockRange(s)

	// Replace the first <version> tag located outside the <parent> block.
	re := regexp.MustCompile(`<version>\s*[^<\s]+\s*</version>`)
	for _, loc := range re.FindAllStringIndex(s, -1) {
		if outsideParent(loc[0], pStart, pEnd) {
			s = s[:loc[0]] + "<version>" + newVersion + "</version>" + s[loc[1]:]
			return os.WriteFile(pomPath, []byte(s), 0o644)
		}
	}

	// Fallback: insert version after the first <artifactId> outside the parent block.
	artRe := regexp.MustCompile(`<artifactId>\s*[^<\s]+\s*</artifactId>`)
	for _, loc := range artRe.FindAllStringIndex(s, -1) {
		if outsideParent(loc[0], pStart, pEnd) {
			insert := "\n  <version>" + newVersion + "</version>"
			s = s[:loc[1]] + insert + s[loc[1]:]
			return os.WriteFile(pomPath, []byte(s), 0o644)
		}
	}

	// Last resort: insert a version section before </project>.
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

// bumpSemver increments a semver string (major.minor.patch). Any pre-release or
// build suffix (e.g. "-SNAPSHOT", "+build") is stripped before parsing and
// re-appended to the result, so "0.0.1-SNAPSHOT" bumps to "0.0.2-SNAPSHOT".
func bumpSemver(v, kind string) (string, error) {
	orig := strings.TrimSpace(v)
	orig = strings.TrimPrefix(orig, "v")
	core := orig
	suffix := ""
	if i := strings.IndexAny(core, "-+"); i >= 0 {
		suffix = core[i:]
		core = core[:i]
	}
	parts := strings.Split(core, ".")
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
	return fmt.Sprintf("%d.%d.%d", major, minor, patch) + suffix, nil
}
