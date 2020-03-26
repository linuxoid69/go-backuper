package backupdb

import (
	"fmt"
	bf "github.com/linuxoid69/go-backuper/backupfiles"
	"github.com/linuxoid69/go-backuper/config"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	// PGDumpCmd bin execute
	PGDumpCmd = "pg_dump"
)

// CreateBackupPostgresDB - create postgres backup
func CreateBackupPostgresDB(cfg *config.Config, db string) error {
	pgOptions := []string{
		fmt.Sprintf("-h%v", cfg.Database.Host),
		fmt.Sprintf("-p%v", cfg.Database.Port),
		fmt.Sprintf("-U%v", cfg.Database.User),
		fmt.Sprintf("-f%v.sql", cfg.BackupStoragePath+"/"+db),
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
func CreateAllPostgresDB(cfg *config.Config) (SQLList []string, err error) {
	for _, i := range cfg.Database.DBnames {
		log.Printf("Create backup for %v", i)

		err := CreateBackupPostgresDB(cfg, i)

		if err != nil {
			log.Printf("Error create backup db '%v' -  %v", i, err)
			return []string{}, err
		}
		SQLList = append(SQLList, cfg.BackupStoragePath+"/"+i+".sql")
		log.Printf("Creating backup file SQL for db '%v'  - success", i)

	}
	return SQLList, nil
}

// CreateArchiveDB create zip archives for db
func CreateArchiveDB(SQLList []string) error {
	for _, i := range SQLList {
		err := bf.CreateArchive([]string{i}, i+".zip")

		if err != nil {
			log.Printf("Creating ZIP archive for '%v' - by path %v", i, i+".zip")
			return err
		}

		log.Printf("Creating ZIP archive for '%v' - by path %v", i, i+".zip")

		err = os.Remove(i)

		if err != nil {
			log.Printf("Can't delete file %v. Error: %v", i, err)
			return err
		}

		log.Printf("File successful deleted '%v'", i)
	}
	return nil
}
