package main

import (
	"fmt"
	"log"

	"github.com/linuxoid69/go-backuper/cmd"
	"github.com/linuxoid69/go-backuper/config"
)

var (
	cfg     *config.Config
	err     error
	version string
	commit  string
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	fmt.Printf("Version: %s commit: %s\n", version, commit)
	cmd.Run()
}
