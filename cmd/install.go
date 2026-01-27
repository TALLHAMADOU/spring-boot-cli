package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var installName string
var installPackage string

var installCmd = &cobra.Command{
	Use:   "install:project [build-tool]",
	Short: "Create a new Spring Boot project (maven|gradle)",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("requires one argument: gradle or maven")
		}
		if args[0] != "gradle" && args[0] != "maven" {
			return fmt.Errorf("build-tool must be 'gradle' or 'maven'")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		tool := args[0]
		name := installName
		if strings.TrimSpace(name) == "" {
			name = "my-spring-project"
		}

		base := filepath.Clean(name)
		// package (group) to use for generated sources
		pkg := installPackage
		if strings.TrimSpace(pkg) == "" {
			pkg = "com.example"
		}
		pkgPath := filepath.Join(strings.Split(pkg, ".")...)
		srcDir := filepath.Join(base, "src", "main", "java", pkgPath)

		if err := os.MkdirAll(srcDir, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "failed to create project directories: %v\n", err)
			return
		}

		// Application.java
		appPath := filepath.Join(srcDir, "Application.java")
		appContent := fmt.Sprintf("package %s;\n\nimport org.springframework.boot.SpringApplication;\nimport org.springframework.boot.autoconfigure.SpringBootApplication;\n\n@SpringBootApplication\npublic class Application {\n    public static void main(String[] args) {\n        SpringApplication.run(Application.class, args);\n    }\n}\n", pkg)
		if err := os.WriteFile(appPath, []byte(appContent), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write Application.java: %v\n", err)
			return
		}

		if tool == "maven" {
			pom := `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
        xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
        xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 https://maven.apache.org/xsd/maven-4.0.0.xsd">
  <modelVersion>4.0.0</modelVersion>
  <parent>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-parent</artifactId>
    <version>3.2.0</version>
    <relativePath/> <!-- lookup parent from repository -->
  </parent>
	<groupId>` + pkg + `</groupId>
  <artifactId>` + name + `</artifactId>
  <version>0.0.1-SNAPSHOT</version>
  <properties>
    <java.version>17</java.version>
  </properties>
  <dependencies>
    <dependency>
      <groupId>org.springframework.boot</groupId>
      <artifactId>spring-boot-starter-web</artifactId>
    </dependency>
    <dependency>
      <groupId>org.springframework.boot</groupId>
      <artifactId>spring-boot-starter-test</artifactId>
      <scope>test</scope>
    </dependency>
  </dependencies>
  <build>
    <plugins>
      <plugin>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-maven-plugin</artifactId>
      </plugin>
    </plugins>
  </build>
</project>
`
			if err := os.WriteFile(filepath.Join(base, "pom.xml"), []byte(pom), 0o644); err != nil {
				fmt.Fprintf(os.Stderr, "failed to write pom.xml: %v\n", err)
				return
			}

			// lightweight mvnw wrappers (shell and windows cmd)
			mvnw := "#!/bin/sh\nif command -v mvn >/dev/null 2>&1; then\n  mvn \"$@\"\nelse\n  echo \"Maven not found. Install Maven or use a machine with Maven.\"\n  exit 1\nfi\n"
			mvnwCmd := "@echo off\nwhere mvn >nul 2>&1 || (echo Maven not found. Please install Maven. & exit /b 1)\nmvn %*\n"
			_ = os.WriteFile(filepath.Join(base, "mvnw"), []byte(mvnw), 0o755)
			_ = os.WriteFile(filepath.Join(base, "mvnw.cmd"), []byte(mvnwCmd), 0o644)

			fmt.Printf("Created Maven project '%s'.\nNext steps:\n  cd %s && mvn spring-boot:run\n", name, name)
			return
		}

		// Gradle
		build := `plugins {
    id 'java'
    id 'org.springframework.boot' version '3.2.0'
    id 'io.spring.dependency-management' version '1.1.0'
}

group = '` + pkg + `'
version = '0.0.1-SNAPSHOT'
sourceCompatibility = '17'

repositories {
    mavenCentral()
}

dependencies {
    implementation 'org.springframework.boot:spring-boot-starter-web'
    testImplementation 'org.springframework.boot:spring-boot-starter-test'
}

tasks.named('test') {
    useJUnitPlatform()
}
`
		if err := os.WriteFile(filepath.Join(base, "build.gradle"), []byte(build), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write build.gradle: %v\n", err)
			return
		}
		settings := "rootProject.name = '" + name + "'\n"
		if err := os.WriteFile(filepath.Join(base, "settings.gradle"), []byte(settings), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write settings.gradle: %v\n", err)
			return
		}

		// lightweight gradle wrappers
		gradlew := "#!/bin/sh\nif command -v gradle >/dev/null 2>&1; then\n  gradle \"$@\"\nelse\n  echo \"Gradle not found. Install Gradle or use the full wrapper.\"\n  exit 1\nfi\n"
		gradlewCmd := "@echo off\nwhere gradle >nul 2>&1 || (echo Gradle not found. Please install Gradle. & exit /b 1)\ngradle %*\n"
		_ = os.WriteFile(filepath.Join(base, "gradlew"), []byte(gradlew), 0o755)
		_ = os.WriteFile(filepath.Join(base, "gradlew.bat"), []byte(gradlewCmd), 0o644)

		fmt.Printf("Created Gradle project '%s'.\nNext steps:\n  cd %s && ./gradlew bootRun\n", name, name)
	},
}

func init() {
	installCmd.Flags().StringVarP(&installName, "name", "n", "my-spring-project", "project name")
	installCmd.Flags().StringVarP(&installPackage, "package", "p", "", "base package (group id)")
	rootCmd.AddCommand(installCmd)
}
