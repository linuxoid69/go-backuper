package cron

import (
	"time"

	"github.com/linuxoid69/go-backuper/config"
	bdb "github.com/linuxoid69/go-backuper/internal/backupdb"
	bf "github.com/linuxoid69/go-backuper/internal/backupfiles"
	en "github.com/linuxoid69/go-backuper/internal/encryption"
	"github.com/linuxoid69/go-backuper/internal/retention"
	cron "github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

// TasksCron - function run tasks of Cron.
func TasksCron(cfg *config.Config) {
	location, err := time.LoadLocation(cfg.TimeZone)
	if err != nil {
		logrus.Error(err)
	}

	c := cron.New()

	cron.WithLocation(location)

	logrus.Infof("Time zone: %s", location)

	c.AddFunc(cfg.CronFiles, func() {
		backupFiles := bf.CreateArchives(cfg)
		en.RunEncrypt(backupFiles, cfg)
	})

	c.AddFunc(cfg.CronDB, func() {
		SQLFiles, err := bdb.CreateAllPostgresDB(cfg)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"message": "Error getting list of SQL files",
				"value":   err}).Error()
		}

		backupDB, err := bdb.CreateArchiveDB(SQLFiles)
		if err != nil {
			logrus.Error(err)
		}

		en.RunEncrypt(backupDB, cfg)
	})

	c.AddFunc("0 */1 * * *", func() {
		if err := retention.RetentionPolicy(cfg); err != nil {
			logrus.Error(err)
		}
	})

	c.Start()
}
