package cmd

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/inconshreveable/go-update"
	"github.com/spf13/cobra"
)

// updateRepo is the GitHub repository whose releases the binary updates from.
const updateRepo = "TALLHAMADOU/spring-boot-cli"

type ghAsset struct {
	Name string `json:"name"`
	URL  string `json:"browser_download_url"`
}

type ghRelease struct {
	TagName string    `json:"tag_name"`
	Assets  []ghAsset `json:"assets"`
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Mettre à jour spring-cli vers la dernière version",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runUpdate(context.Background()); err != nil {
			Error("%v", err)
		}
	},
}

func runUpdate(ctx context.Context) error {
	Info("Recherche de mises à jour...")

	rel, err := fetchLatestRelease(ctx, updateRepo)
	if err != nil {
		return err
	}

	latestVer := strings.TrimPrefix(rel.TagName, "v")
	currentVer := strings.TrimPrefix(Version, "v")
	if currentVer != "dev" && currentVer == latestVer {
		Success("Vous utilisez déjà la dernière version : v%s", currentVer)
		return nil
	}

	asset, ok := selectAsset(rel.Assets, runtime.GOOS, runtime.GOARCH)
	if !ok {
		return fmt.Errorf("aucune archive disponible pour %s/%s dans la release %s", runtime.GOOS, runtime.GOARCH, rel.TagName)
	}

	Info("Nouvelle version %s. Téléchargement de %s...", rel.TagName, asset.Name)
	bin, err := downloadBinary(ctx, asset.URL)
	if err != nil {
		return fmt.Errorf("téléchargement de la mise à jour: %w", err)
	}

	if err := update.Apply(bytes.NewReader(bin), update.Options{}); err != nil {
		return fmt.Errorf("application de la mise à jour: %w", err)
	}
	Success("spring-cli a été mis à jour vers %s !", rel.TagName)
	return nil
}

// fetchLatestRelease returns the latest GitHub release for repo (owner/name).
func fetchLatestRelease(ctx context.Context, repo string) (*ghRelease, error) {
	url := "https://api.github.com/repos/" + repo + "/releases/latest"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := (&http.Client{Timeout: 30 * time.Second}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, fmt.Errorf("aucune release publiée pour %s", repo)
	default:
		return nil, fmt.Errorf("réponse GitHub inattendue: %s", resp.Status)
	}

	var rel ghRelease
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, err
	}
	return &rel, nil
}

// selectAsset picks the release archive matching the given OS/architecture,
// following the goreleaser naming convention (…_<os>_<arch>.tar.gz).
func selectAsset(assets []ghAsset, goos, goarch string) (ghAsset, bool) {
	for _, a := range assets {
		n := strings.ToLower(a.Name)
		if strings.HasSuffix(n, ".tar.gz") && strings.Contains(n, goos) && strings.Contains(n, goarch) {
			return a, true
		}
	}
	return ghAsset{}, false
}

// downloadBinary fetches a .tar.gz release archive and returns the spring-cli
// executable it contains.
func downloadBinary(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := (&http.Client{Timeout: 5 * time.Minute}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("téléchargement %s: %s", url, resp.Status)
	}

	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		base := path.Base(hdr.Name)
		if base == "spring-cli" || base == "spring-cli.exe" {
			return io.ReadAll(tr)
		}
	}
	return nil, fmt.Errorf("binaire spring-cli introuvable dans l'archive")
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
