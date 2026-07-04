/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"os"

	"github.com/TALLHAMADOU/spring-boot-cli/cmd"
)

func main() {
	// support `spring-cli --version` and `spring-cli -v` directly without invoking Cobra
	if len(os.Args) == 2 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("spring-cli version %s\n", cmd.Version)
		return
	}

	cmd.Execute()
}
