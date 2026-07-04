package cmd

import "testing"

func TestSelectAsset(t *testing.T) {
	assets := []ghAsset{
		{Name: "spring-cli_1.2.3_linux_amd64.tar.gz", URL: "u1"},
		{Name: "spring-cli_1.2.3_linux_arm64.tar.gz", URL: "u2"},
		{Name: "spring-cli_1.2.3_darwin_arm64.tar.gz", URL: "u3"},
		{Name: "spring-cli_1.2.3_windows_amd64.tar.gz", URL: "u4"},
		{Name: "checksums.txt", URL: "u5"},
	}

	cases := []struct {
		goos, goarch string
		wantURL      string
		wantOK       bool
	}{
		{"linux", "amd64", "u1", true},
		{"linux", "arm64", "u2", true},
		{"darwin", "arm64", "u3", true},
		{"windows", "amd64", "u4", true},
		{"freebsd", "amd64", "", false},
		{"linux", "386", "", false},
	}

	for _, c := range cases {
		got, ok := selectAsset(assets, c.goos, c.goarch)
		if ok != c.wantOK || got.URL != c.wantURL {
			t.Errorf("selectAsset(%s/%s) = (%q,%v), want (%q,%v)",
				c.goos, c.goarch, got.URL, ok, c.wantURL, c.wantOK)
		}
	}
}
