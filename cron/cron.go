package cron

import (
	"github.com/robfig/cron"
	"time"
	bdb "github.com/linuxoid69/go-backuper/backupdb"
	bf "github.com/linuxoid69/go-backuper/backupfiles"
	"log"
	en "github.com/linuxoid69/go-backuper/encryption"
	"github.com/linuxoid69/go-backuper/config"
)

// TasksCron - function run tasks of Cron
func TasksCron(cfg *config.Config)  {
	l, _ := time.LoadLocation(cfg.TimeZone)
	c := cron.New()
	cron.NewWithLocation(l)
	log.Printf("Time zone: %s", l)

	c.AddFunc(cfg.CronFiles, func() {
		backupFiles := bf.CreateArchives(cfg)
		en.RunEncrypt(backupFiles, cfg)
	})

	c.AddFunc(cfg.CronDB, func() {
		SQLFiles, err := bdb.CreateAllPostgresDB(cfg)
		if err != nil {
			log.Printf("Error getting list of SQL files %v", err)
		}
		backupDB, _ := bdb.CreateArchiveDB(SQLFiles)
		en.RunEncrypt(backupDB, cfg)
	})

	c.Start()
}
