package backupdb

import (
	"fmt"
	bf "github.com/linuxoid69/go-backuper/backupfiles"
	"github.com/linuxoid69/go-backuper/config"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	// PGDumpCmd bin execute
	PGDumpCmd = "pg_dump"
)

// CreateBackupPostgresDB - create postgres backup
func CreateBackupPostgresDB(cfg *config.Config, db string) error {
	t := strings.Replace(time.Now().Format("02-01-2006"), "-", "_", 2)
	pgOptions := []string{
		fmt.Sprintf("-h%v", cfg.Database.Host),
		fmt.Sprintf("-p%v", cfg.Database.Port),
		fmt.Sprintf("-U%v", cfg.Database.User),
		fmt.Sprintf("-f%v_%v.sql", cfg.BackupStoragePath+"/"+db, t),
		fmt.Sprintf("-d%v", db),
		fmt.Sprintf("-w%v", ""),
	}

	if cfg.Database.Options != "" {
		pgOptions = append(pgOptions, strings.Split(cfg.Database.Options, " ")...)
	}

	os.Setenv("PGPASSWORD", cfg.Database.Password)

	cmd := exec.Command(PGDumpCmd, pgOptions...)
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("Error %v + %v", err, string(out))
		return err
	}

	if len(out) > 0 {
		log.Println(string(out))
	}

	return nil
}

// CreateAllPostgresDB - create all DB from config
func CreateAllPostgresDB(cfg *config.Config) (dumpList []string, err error) {
	t := strings.Replace(time.Now().Format("02-01-2006"), "-", "_", 2)
	for _, i := range cfg.Database.DBnames {
		log.Printf("Create backup for %v", i)

		err := CreateBackupPostgresDB(cfg, i)

		if err != nil {
			log.Printf("Error create backup db '%v' -  %v", i, err)
			return []string{}, err
		}
		dumpFilePath := fmt.Sprintf("%v/%v_%v.sql", cfg.BackupStoragePath, i, t)
		dumpList = append(dumpList, dumpFilePath)
		log.Printf("Creating backup file SQL for db '%v'  - success", i)

	}
	return dumpList, nil
}

// CreateArchiveDB create zip archives for db
func CreateArchiveDB(dumpList []string) (backupFiles []string, err error) {
	for _, i := range dumpList {
		backupFile := fmt.Sprintf("%v.zip", i)
		err = bf.CreateArchive([]string{i}, backupFile)

		if err != nil {
			log.Printf("Creating ZIP archive for '%v' - by path %v", i, i+".zip")
			return
		}

		log.Printf("Creating ZIP archive for '%v' - by path %v", i, i+".zip")

		err = os.Remove(i)

		if err != nil {
			log.Printf("Can't delete file %v. Error: %v", i, err)
			return
		}

		log.Printf("File successful deleted '%v'", i)
		backupFiles = append(backupFiles, backupFile)
	}
	return backupFiles, nil
}
