package main

import (
	"fmt"
	"github.com/crask/mqproxy/internal/version"
	"github.com/crask/mqproxy/server"
	"log"
	"os"
)

func usage(command string) {
	fmt.Printf("Usage of %s:\n", command)
	fmt.Println("\t-c config_file:\tRunning mqproxy with config")
	fmt.Println("\t-v:\tPrint version")
}

func checkArgs(args []string) bool {
	return !(len(os.Args) < 2 ||
		(os.Args[1] != "-c" && os.Args[1] != "-v") ||
		(os.Args[1] == "-c" && len(os.Args) < 3))
}

func main() {
	if !checkArgs(os.Args) {
		usage(os.Args[0])
		os.Exit(1)
	}

	if os.Args[1] == "-v" {
		fmt.Println(version.String("mqproxy"))
		os.Exit(0)
	}

	cfg, err := server.NewProxyConfig(os.Args[2])
	if err != nil {
		log.Fatalf("parse config error, %v", err)
	}

	if err = server.Startable(cfg); err != nil {
		log.Fatalf("server startable error, %v", err)
	}
}
