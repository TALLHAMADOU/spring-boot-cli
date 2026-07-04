package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBumpSemver(t *testing.T) {
	v, err := bumpSemver("1.2.3", "patch")
	if err != nil || v != "1.2.4" {
		t.Fatalf("expected 1.2.4 got %s err %v", v, err)
	}
	v, _ = bumpSemver("1.2.3", "minor")
	if v != "1.3.0" {
		t.Fatalf("expected 1.3.0 got %s", v)
	}
	v, _ = bumpSemver("1.2.3", "major")
	if v != "2.0.0" {
		t.Fatalf("expected 2.0.0 got %s", v)
	}
	// Pre-release suffix (Spring's default) must be preserved, not rejected.
	v, err = bumpSemver("0.0.1-SNAPSHOT", "patch")
	if err != nil || v != "0.0.2-SNAPSHOT" {
		t.Fatalf("expected 0.0.2-SNAPSHOT got %s err %v", v, err)
	}
	v, _ = bumpSemver("1.4.2-SNAPSHOT", "minor")
	if v != "1.5.0-SNAPSHOT" {
		t.Fatalf("expected 1.5.0-SNAPSHOT got %s", v)
	}
}

func TestDetectVersionInPOMandSet(t *testing.T) {
	tmp := t.TempDir()
	pom := `<?xml version="1.0" encoding="UTF-8"?>
<project>
  <modelVersion>4.0.0</modelVersion>
  <groupId>com.example</groupId>
  <artifactId>testproj</artifactId>
  <version>0.0.1-SNAPSHOT</version>
</project>`
	p := filepath.Join(tmp, "pom.xml")
	if err := os.WriteFile(p, []byte(pom), 0o644); err != nil {
		t.Fatal(err)
	}
	v, err := detectVersionInPOM(tmp)
	if err != nil || v != "0.0.1-SNAPSHOT" {
		t.Fatalf("expected 0.0.1-SNAPSHOT got %s err %v", v, err)
	}
	if err := setVersionInPOM(tmp, "0.1.0"); err != nil {
		t.Fatal(err)
	}
	v2, err := detectVersionInPOM(tmp)
	if err != nil || v2 != "0.1.0" {
		t.Fatalf("expected 0.1.0 got %s err %v", v2, err)
	}
}

func TestSetVersionInPOMWithParent(t *testing.T) {
	tmp := t.TempDir()
	pom := `<?xml version="1.0" encoding="UTF-8"?>
<project>
  <modelVersion>4.0.0</modelVersion>
  <parent>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-parent</artifactId>
    <version>3.2.0</version>
  </parent>
  <groupId>com.example</groupId>
  <artifactId>demo</artifactId>
  <version>0.0.1-SNAPSHOT</version>
</project>`
	p := filepath.Join(tmp, "pom.xml")
	if err := os.WriteFile(p, []byte(pom), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := setVersionInPOM(tmp, "0.2.0"); err != nil {
		t.Fatal(err)
	}
	out, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	// Parent version must be preserved.
	if !strings.Contains(s, "<version>3.2.0</version>") {
		t.Fatalf("parent version was modified:\n%s", s)
	}
	// Project version must be updated.
	if !strings.Contains(s, "<version>0.2.0</version>") {
		t.Fatalf("project version not set to 0.2.0:\n%s", s)
	}
	if strings.Contains(s, "0.0.1-SNAPSHOT") {
		t.Fatalf("old project version still present:\n%s", s)
	}
	// detectVersionInPOM must report the project version, not the parent one.
	v, err := detectVersionInPOM(tmp)
	if err != nil || v != "0.2.0" {
		t.Fatalf("expected detected 0.2.0 got %s err %v", v, err)
	}
}

func TestDetectVersionInGradleAndSet(t *testing.T) {
	tmp := t.TempDir()
	grad := "version = '0.0.1-SNAPSHOT'\n"
	if err := os.WriteFile(filepath.Join(tmp, "build.gradle"), []byte(grad), 0o644); err != nil {
		t.Fatal(err)
	}
	v, err := detectVersionInGradle(tmp)
	if err != nil || v != "0.0.1-SNAPSHOT" {
		t.Fatalf("expected gradle version found got %s err %v", v, err)
	}
	if err := setVersionInGradle(tmp, "0.1.0"); err != nil {
		t.Fatal(err)
	}
	v2, err := detectVersionInGradle(tmp)
	if err != nil || v2 != "0.1.0" {
		t.Fatalf("expected gradle updated to 0.1.0 got %s err %v", v2, err)
	}
}
