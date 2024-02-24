package main

import "github.com/lucafroeschke/go-package-server/server"

func main() {
	err := server.Start()
	if err != nil {
		panic(err)
	}
}
