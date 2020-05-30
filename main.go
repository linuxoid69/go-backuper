package main

import (
	"github.com/linuxoid69/go-backuper/cmd"
	"github.com/linuxoid69/go-backuper/config"
	"log"
)

var (
	cfg *config.Config
	err error
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	cmd.Run()
}
