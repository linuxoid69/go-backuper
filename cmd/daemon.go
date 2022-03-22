package cmd

import (
	"os"
	"time"

	"github.com/linuxoid69/go-backuper/config"
	"github.com/linuxoid69/go-backuper/internal/cron"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals
var configPath *string

// daemonCmd represents the daemon command.
//nolint:gochecknoglobals
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run application in daemon mode",
	Long:  "Run application in daemon mode",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(*configPath)
		if err != nil {
			logrus.Error(err)
			os.Exit(1)
		}

		cron.TasksCron(cfg)

		for {
			//nolint:gomnd
			time.Sleep(time.Millisecond * 10)
		}
	},
}

//nolint:gochecknoinits
func init() {
	rootCmd.AddCommand(daemonCmd)
	configPath = daemonCmd.PersistentFlags().StringP("config", "c", "", "config")
}
