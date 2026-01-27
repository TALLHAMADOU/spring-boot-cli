package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var (
	versionBump   string
	versionSet    string
	versionAuto   bool
	versionCommit bool
	versionTag    bool
	versionPush   bool
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Inspect and modify project version (pom.xml / build.gradle)",
	Run: func(cmd *cobra.Command, args []string) {
		if !isSpringProject(".") {
			fmt.Fprintln(os.Stderr, "Erreur: Lancez cette commande dans un projet Spring Boot (pom.xml ou build.gradle requis)")
			os.Exit(1)
		}

		cur, kind, err := detectVersionInBuildFiles(".")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to detect version: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Current version (%s): %s\n", kind, cur)

		var newVer string
		if versionSet != "" {
			newVer = versionSet
		} else if versionBump != "" {
			nv, err := bumpSemver(cur, versionBump)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to bump version: %v\n", err)
				os.Exit(1)
			}
			newVer = nv
		} else if versionAuto {
			// get last git tag
			out, err := exec.Command("git", "describe", "--tags", "--abbrev=0").Output()
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to get last git tag: %v\n", err)
				os.Exit(1)
			}
			last := strings.TrimSpace(string(out))
			last = strings.TrimPrefix(last, "v")
			nv, err := bumpSemver(last, "patch")
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to bump auto version: %v\n", err)
				os.Exit(1)
			}
			newVer = nv
			fmt.Printf("Auto new version (from %s): %s\n", last, newVer)
			if !versionCommit && !versionTag {
				fmt.Println("--auto used without --commit/--tag: no changes applied (use --commit to apply)")
				return
			}
		} else {
			// nothing to do
			return
		}

		// apply change
		switch kind {
		case "maven":
			if err := setVersionInPOM(".", newVer); err != nil {
				fmt.Fprintf(os.Stderr, "failed to set version in pom.xml: %v\n", err)
				os.Exit(1)
			}
		case "gradle":
			if err := setVersionInGradle(".", newVer); err != nil {
				fmt.Fprintf(os.Stderr, "failed to set version in build.gradle: %v\n", err)
				os.Exit(1)
			}
		default:
			fmt.Fprintf(os.Stderr, "unsupported build file kind: %s\n", kind)
			os.Exit(1)
		}
		fmt.Printf("Updated version to %s\n", newVer)

		if versionCommit {
			_ = exec.Command("git", "add", "pom.xml", "build.gradle").Run()
			msg := "chore: bump version to " + newVer
			if err := exec.Command("git", "commit", "-m", msg).Run(); err != nil {
				fmt.Fprintf(os.Stderr, "git commit failed: %v\n", err)
			} else {
				fmt.Println("Committed version change")
			}
		}

		if versionTag {
			tag := "v" + newVer
			if err := exec.Command("git", "tag", tag).Run(); err != nil {
				fmt.Fprintf(os.Stderr, "git tag failed: %v\n", err)
			} else {
				fmt.Printf("Created tag %s\n", tag)
			}
		}

		if versionPush {
			if versionCommit {
				if err := exec.Command("git", "push").Run(); err != nil {
					fmt.Fprintf(os.Stderr, "git push failed: %v\n", err)
				} else {
					fmt.Println("Pushed commits")
				}
			}
			if versionTag {
				if err := exec.Command("git", "push", "--tags").Run(); err != nil {
					fmt.Fprintf(os.Stderr, "git push tags failed: %v\n", err)
				} else {
					fmt.Println("Pushed tags")
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().StringVar(&versionBump, "bump", "", "bump kind: patch|minor|major")
	versionCmd.Flags().StringVar(&versionSet, "set", "", "set explicit version x.y.z")
	versionCmd.Flags().BoolVar(&versionAuto, "auto", false, "auto bump from last git tag (increment patch)")
	versionCmd.Flags().BoolVar(&versionCommit, "commit", false, "git add & commit the version change")
	versionCmd.Flags().BoolVar(&versionTag, "tag", false, "create git tag vX.Y.Z")
	versionCmd.Flags().BoolVar(&versionPush, "push", false, "push commits and/or tags")
}
