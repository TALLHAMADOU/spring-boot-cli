package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var dockerDB string
var dockerPort int
var dockerJDK int

var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Generate Dockerfile and docker-compose.yml",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireSpringProject(); err != nil {
			return err
		}

		// Generate Dockerfile
		dockerfileContent, err := renderTemplate("dockerfile", struct{ JDK int }{JDK: dockerJDK})
		if err != nil {
			return fmt.Errorf("rendu du template Dockerfile: %w", err)
		}
		if err := os.WriteFile("Dockerfile", []byte(dockerfileContent), 0o644); err != nil {
			return fmt.Errorf("écriture du Dockerfile: %w", err)
		}
		Success("Created Dockerfile\n")

		// Generate docker-compose.yml
		composeContent, err := renderTemplate("docker_compose", struct {
			DB   string
			Port int
		}{DB: dockerDB, Port: dockerPort})
		if err != nil {
			return fmt.Errorf("rendu du template docker_compose: %w", err)
		}
		if err := os.WriteFile("docker-compose.yml", []byte(composeContent), 0o644); err != nil {
			return fmt.Errorf("écriture du docker-compose.yml: %w", err)
		}
		Success("Created docker-compose.yml\n")

		return nil
	},
}

func init() {
	dockerCmd.Flags().StringVar(&dockerDB, "db", "postgres", "Database to include in docker-compose (postgres, mysql)")
	dockerCmd.Flags().IntVar(&dockerPort, "port", 8080, "Port exposed by the Spring Boot application")
	dockerCmd.Flags().IntVar(&dockerJDK, "jdk", 17, "JDK version to use in the Dockerfile")
	makeCmd.AddCommand(dockerCmd)
}
