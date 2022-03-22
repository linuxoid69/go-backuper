package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config - struct of config
type Config struct {
	BackupStoragePath     string                           `yaml:"backup_storage_path"`
	BackupRetentionDays   int                              `yaml:"backup_retention_days"`
	NameZipFile           string                           `yaml:"name_zip_file"`
	BackupSourceFilesPath map[string][][]map[string]string `yaml:"backup_source_files_path"`
	CronFiles             string                           `yaml:"cron_files"`
	CronDB                string                           `yaml:"cron_db"`
	TimeZone              string                           `yaml:"time_zone"`
	EncryptBackup         bool                             `yaml:"encrypt_backup"`
	EncryptPassword       string                           `yaml:"encrypt_password"`
	Database              struct {
		Host     string   `yaml:"host"`
		Port     int      `yaml:"port"`
		User     string   `yaml:"user"`
		Password string   `yaml:"password"`
		Options  string   `yaml:"options"`
		DBnames  []string `yaml:",flow"`
	} `yaml:"database"`
}

// readConfig - read config file
func readConfig(path string) ([]byte, error) {
	dat, err := ioutil.ReadFile(path)

	return dat, err
}

// LoadConfig -load config from config.yaml.
func LoadConfig(configPath string) (*Config, error) {
	conf := Config{}

	dat, err := readConfig(configPath)
	if err != nil {
		return &conf, err
	}

	err = yaml.Unmarshal([]byte(dat), &conf)
	if err != nil {
		return &conf, err
	}

	return &conf, nil
}
