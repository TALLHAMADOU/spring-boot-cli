package cmd

import (
	"context"
	"os"
	"strings"

	"github.com/creativeprojects/go-selfupdate"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Mettre à jour spring-cli vers la dernière version",
	Run: func(cmd *cobra.Command, args []string) {
		Info("Recherche de mises à jour...")

		// Parse the repository slug
		repo := "TALLHAMADOU/spring-boot-cli"

		latest, found, err := selfupdate.DetectLatest(context.Background(), selfupdate.ParseSlug(repo))
		if err != nil {
			Error("Erreur lors de la recherche de mises à jour: %v", err)
			return
		}

		if !found {
			Warning("Aucune release trouvée pour le dépôt %s.", repo)
			return
		}

		// Normalize versions for comparison
		currentVer := strings.TrimPrefix(Version, "v")
		latestVer := strings.TrimPrefix(latest.Version(), "v")

		if currentVer != "dev" && currentVer == latestVer {
			Success("Vous utilisez déjà la dernière version : v%s", currentVer)
			return
		}

		Info("Nouvelle version trouvée : %s. Téléchargement et mise à jour en cours...", latest.Version())

		exe, err := os.Executable()
		if err != nil {
			Error("Impossible de localiser le chemin de l'exécutable: %v", err)
			return
		}

		if err := selfupdate.UpdateTo(context.Background(), latest.AssetURL, latest.AssetName, exe); err != nil {
			Error("Échec de la mise à jour : %v", err)
			return
		}

		Success("spring-cli a été mis à jour avec succès vers la version %s !", latest.Version())
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
