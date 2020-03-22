package main

import (
	"github.com/robfig/cron"
	bdb "github.com/linuxoid69/go-backuper/backupdb"
	bf "github.com/linuxoid69/go-backuper/backupfiles"
	"github.com/linuxoid69/go-backuper/config"
	"log"
	"time"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Printf("Start application")
}

func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("Can't load config %v: ", err)
	}
	l,_ := time.LoadLocation(cfg.TimeZone)
	c := cron.New(
		cron.WithLocation(l),
	)
	log.Println(l)
	c.AddFunc(cfg.CronFiles, func() {
		bf.CreateArchives(cfg)
	})

	c.AddFunc(cfg.CronDB, func() {
		SQLFiles, err := bdb.CreateAllPostgresDB(cfg)
		if err != nil {
			log.Printf("Error getting list of SQL files %v", err)
		}
		bdb.CreateArchiveDB(SQLFiles)
	})

	c.Start()

	for {
		time.Sleep(time.Millisecond * 10)
	}
}
