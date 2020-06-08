package retention

import (
	"github.com/linuxoid69/go-backuper/config"
	"github.com/linuxoid69/go-backuper/helpers"
	"os"
)

// RetentionPolicy - delete old backups
func RetentionPolicy(cfg *config.Config) error {
	// actions
	files, err := helpers.Find(cfg.BackupStoragePath, "file", cfg.BackupRetentionDays)
	if err != nil {
		return err
	}
	for _, file := range files {
		err = os.Remove(file)
		if err != nil {
			return err
		}
	}
	return nil
}
