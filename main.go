package main

import (
	"log"
	"os"

	"github.com/lucafroeschke/go-package-server/cmd"
	"github.com/lucafroeschke/go-package-server/server"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "build":
			cmd.RunBuildCommand()
		case "config":
			cmd.RunConfigCommand()
		}
	} else {
		err := server.Start()
		if err != nil {
			log.Fatal(err)
		}
	}
}
